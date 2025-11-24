package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Pour gérer gorm.ErrRecordNotFound
)

// ClickEventsChan est le channel global utilisé pour envoyer les événements
var ClickEventsChan chan<- models.ClickEvent

// TODO Créer une variable ClickEventsChannel qui est un chan de type ClickEvent
// ClickEventsChannel est le channel global (ou injecté) utilisé pour envoyer les événements de clic
// aux workers asynchrones. Il est bufferisé pour ne pas bloquer les requêtes de redirection.
// PHASE 3 : Pas utilisé pour l'instant (sans async)

// SetupRoutes configure toutes les routes de l'API Gin et injecte les dépendances nécessaires
func SetupRoutes(router *gin.Engine, linkService *services.LinkService, clickService *services.ClickService) {
	// PHASE 3 : Pas de channel pour l'instant (sans async)

	// DONE : Route de Health Check , /health
	router.GET("/health", HealthCheckHandler)

	// DONE : Routes de l'API
	// Doivent être au format /api/v1/
	router.POST("/api/v1/links", CreateShortLinkHandler(linkService))
	router.GET("/api/v1/links/:shortCode", GetLinkInfoHandler(linkService))
	router.GET("/api/v1/links/:shortCode/stats", GetLinkStatsHandler(linkService))

	// Route de Redirection (au niveau racine pour les short codes)
	// IMPORTANT: Doit être APRÈS les routes /api/v1/ pour éviter les conflits
	router.GET("/:shortCode", RedirectHandler(linkService, clickService))
}

// HealthCheckHandler gère la route /health pour vérifier l'état du service.
func HealthCheckHandler(c *gin.Context) {
	// DONE  Retourner simplement du JSON avec un StatusOK, {"status": "ok"}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	c.Writer.Write([]byte("\n"))
}

// CreateLinkRequest représente le corps de la requête JSON pour la création d'un lien.
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"` // 'binding:required' pour validation, 'url' pour format URL
}

// CreateShortLinkHandler gère la création d'une URL courte.
func CreateShortLinkHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		// DONE : Tente de lier le JSON de la requête à la structure CreateLinkRequest.
		// Gin gère la validation 'binding'.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		// DONE: Appeler le LinkService (CreateLink) pour créer le nouveau lien.
		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			log.Printf("Error creating link: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short link, existed url"})
			c.Writer.Write([]byte("\n"))
			return
		}

		// Retourne le code court et l'URL longue dans la réponse JSON.
		// DONE Choisir le bon code HTTP
		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.Shortcode,
			"long_url":       link.LongURL,
			"full_short_url": "http://localhost:8080/" + link.Shortcode,
		})
		c.Writer.Write([]byte("\n"))
	}
}

// RedirectHandler gère la redirection d'une URL courte vers l'URL longue.
func RedirectHandler(linkService *services.LinkService, clickService *services.ClickService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// DONE Récupère le shortCode de l'URL avec c.Param
		shortCode := c.Param("shortCode")

		// DONE: Récupérer l'URL longue associée au shortCode depuis le linkService (GetLinkByShortCode)
		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			// Si le lien n'est pas trouvé, retourner HTTP 404 Not Found.
			// Utiliser errors.Is et l'erreur Gorm
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
				return
			}
			// Gérer d'autres erreurs potentielles de la base de données ou du service
			log.Printf("Error retrieving link for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Créer un ClickEvent avec les informations pertinentes.
		userAgent := c.Request.UserAgent()
		ipAddress := c.ClientIP()

		// Créer l'objet Click
		click := &models.Click{
			LinkID:    link.ID,
			UserAgent: userAgent,
			IPAddress: ipAddress,
		}

		// DONE : Créer un ClickEvent et l'envoyer dans le channel (async)
		evt := models.ClickEvent{
			LinkID:    link.ID,
			Timestamp: time.Now(),
			UserAgent: userAgent,
			IPAddress: ipAddress,
		}

		if ClickEventsChan != nil {
			select {
			case ClickEventsChan <- evt:
				// envoyé de façon asynchrone
			default:
				log.Printf("WARN: click event dropped for link %d", link.ID)
			}
		} else {
			// Fallback synchrone si le channel n'est pas configuré
			if err := clickService.RecordClick(click); err != nil {
				log.Printf("Error recording click for link %d: %v", link.ID, err)
			}
		}

		// Retourner l'URL en JSON au lieu de rediriger
		c.JSON(http.StatusOK, gin.H{
			"short_code": link.Shortcode,
			"long_url":   link.LongURL,
		})
		c.Writer.Write([]byte("\n"))

		// REDIRECTION HTTP 302 :
		// c.Redirect(http.StatusFound, link.LongURL)
	}
}

// GetLinkInfoHandler gère la récupération des informations d'un lien sans redirection.
func GetLinkInfoHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
				return
			}
			log.Printf("Error retrieving link info for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_code": link.Shortcode,
			"long_url":   link.LongURL,
		})
		c.Writer.Write([]byte("\n"))
	}
}

// GetLinkStatsHandler gère la récupération des statistiques pour un lien spécifique.
func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// DONE Récupère le shortCode de l'URL avec c.Param
		shortCode := c.Param("shortCode")

		// DONE: Appeler le LinkService pour obtenir le lien et le nombre total de clics.
		link, totalClicks, err := linkService.GetLinkStats(shortCode)
		if err != nil {
			// Gérer le cas où le lien n'est pas trouvé.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
				return
			}
			// Gérer d'autres erreurs
			log.Printf("Error retrieving stats for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Retourne les statistiques dans la réponse JSON.
		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.Shortcode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
		c.Writer.Write([]byte("\n"))
	}
}
