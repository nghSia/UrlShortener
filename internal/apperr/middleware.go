package apperr

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorHandler est un middleware Gin qui gère les erreurs de façon centralisée
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Exécuter le handler suivant
		c.Next()

		// Vérifier s'il y a des erreurs
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Si c'est une AppError, l'utiliser directement
			if appErr, ok := err.(*AppError); ok {
				// Logger l'erreur interne si elle existe
				if appErr.InternalErr != nil {
					log.Printf("[ERROR] %s: %v", appErr.Message, appErr.InternalErr)
				} else {
					log.Printf("[ERROR] %s (details: %s)", appErr.Message, appErr.Details)
				}

				// Retourner la réponse JSON
				c.JSON(appErr.Code, appErr.ToJSON())
				return
			}

			// Erreur non gérée, retourner une erreur 500 générique
			log.Printf("[ERROR] Unhandled error: %v", err)
			genericErr := ErrInternalServer("Une erreur inattendue est survenue", err)
			c.JSON(genericErr.Code, genericErr.ToJSON())
		}
	}
}

// HandleError est une fonction helper pour retourner directement une erreur
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		// Logger l'erreur interne si elle existe
		if appErr.InternalErr != nil {
			log.Printf("[ERROR] %s: %v", appErr.Message, appErr.InternalErr)
		} else {
			log.Printf("[ERROR] %s (details: %s)", appErr.Message, appErr.Details)
		}

		c.JSON(appErr.Code, appErr.ToJSON())
		return
	}

	// Erreur non gérée
	log.Printf("[ERROR] Unhandled error: %v", err)
	genericErr := ErrInternalServer("Une erreur inattendue est survenue", err)
	c.JSON(genericErr.Code, genericErr.ToJSON())
}

// AbortWithError arrête l'exécution et retourne une erreur
func AbortWithError(c *gin.Context, err error) {
	HandleError(c, err)
	c.Abort()
}
