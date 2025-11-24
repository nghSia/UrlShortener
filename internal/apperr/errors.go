package apperr

import (
	"fmt"
	"net/http"
)

// AppError représente une erreur applicative avec un code HTTP et un message clair
type AppError struct {
	Code        int    `json:"code"`              // Code HTTP
	Message     string `json:"message"`           // Message utilisateur
	Details     string `json:"details,omitempty"` // Détails techniques (optionnel)
	InternalErr error  `json:"-"`                 // Erreur interne (non exposée)
}

// Error implémente l'interface error
func (e *AppError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.InternalErr)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap permet d'utiliser errors.Is et errors.As
func (e *AppError) Unwrap() error {
	return e.InternalErr
}

// ToJSON retourne la représentation JSON de l'erreur pour l'API
func (e *AppError) ToJSON() map[string]interface{} {
	result := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		},
	}

	if e.Details != "" {
		result["error"].(map[string]interface{})["details"] = e.Details
	}

	return result
}

// NewAppError crée une nouvelle erreur applicative
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:        code,
		Message:     message,
		InternalErr: err,
	}
}

// WithDetails ajoute des détails à l'erreur
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Erreurs 400 - Bad Request

// ErrInvalidRequest retourne une erreur pour une requête invalide
func ErrInvalidRequest(details string, err error) *AppError {
	return &AppError{
		Code:        http.StatusBadRequest,
		Message:     "La requête est invalide",
		Details:     details,
		InternalErr: err,
	}
}

// ErrInvalidURL retourne une erreur pour une URL invalide
func ErrInvalidURL(url string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "L'URL fournie est invalide",
		Details: fmt.Sprintf("URL: %s", url),
	}
}

// ErrMissingField retourne une erreur pour un champ manquant
func ErrMissingField(field string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Champ requis manquant",
		Details: fmt.Sprintf("Le champ '%s' est obligatoire", field),
	}
}

// ErrInvalidShortCode retourne une erreur pour un code court invalide
func ErrInvalidShortCode(shortCode string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Le code court est invalide",
		Details: fmt.Sprintf("Code: %s", shortCode),
	}
}

// Erreurs 404 - Not Found

// ErrLinkNotFound retourne une erreur quand un lien n'est pas trouvé
func ErrLinkNotFound(shortCode string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: "Lien court introuvable",
		Details: fmt.Sprintf("Aucun lien trouvé pour le code: %s", shortCode),
	}
}

// ErrResourceNotFound retourne une erreur pour une ressource introuvable
func ErrResourceNotFound(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: "Ressource introuvable",
		Details: fmt.Sprintf("La ressource '%s' n'existe pas", resource),
	}
}

// Erreurs 409 - Conflict

// ErrLinkAlreadyExists retourne une erreur quand un lien existe déjà
func ErrLinkAlreadyExists(url string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: "Ce lien existe déjà",
		Details: fmt.Sprintf("Un lien court existe déjà pour cette URL: %s", url),
	}
}

// ErrShortCodeAlreadyExists retourne une erreur quand un code court existe déjà
func ErrShortCodeAlreadyExists(shortCode string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: "Ce code court est déjà utilisé",
		Details: fmt.Sprintf("Code: %s", shortCode),
	}
}

// Erreurs 500 - Internal Server Error

// ErrDatabaseOperation retourne une erreur pour un problème de base de données
func ErrDatabaseOperation(operation string, err error) *AppError {
	return &AppError{
		Code:        http.StatusInternalServerError,
		Message:     "Erreur lors de l'opération en base de données",
		Details:     fmt.Sprintf("Opération: %s", operation),
		InternalErr: err,
	}
}

// ErrInternalServer retourne une erreur serveur générique
func ErrInternalServer(details string, err error) *AppError {
	return &AppError{
		Code:        http.StatusInternalServerError,
		Message:     "Une erreur interne est survenue",
		Details:     details,
		InternalErr: err,
	}
}

// ErrFailedToCreateLink retourne une erreur lors de la création d'un lien
func ErrFailedToCreateLink(err error) *AppError {
	return &AppError{
		Code:        http.StatusInternalServerError,
		Message:     "Impossible de créer le lien court",
		Details:     "Une erreur est survenue lors de la création du lien",
		InternalErr: err,
	}
}

// ErrFailedToRecordClick retourne une erreur lors de l'enregistrement d'un clic
func ErrFailedToRecordClick(err error) *AppError {
	return &AppError{
		Code:        http.StatusInternalServerError,
		Message:     "Impossible d'enregistrer le clic",
		Details:     "Le lien fonctionne mais les statistiques n'ont pas pu être enregistrées",
		InternalErr: err,
	}
}

// Erreurs 503 - Service Unavailable

// ErrServiceUnavailable retourne une erreur quand le service est indisponible
func ErrServiceUnavailable(details string) *AppError {
	return &AppError{
		Code:    http.StatusServiceUnavailable,
		Message: "Service temporairement indisponible",
		Details: details,
	}
}

// ErrDatabaseUnavailable retourne une erreur quand la base de données est indisponible
func ErrDatabaseUnavailable(err error) *AppError {
	return &AppError{
		Code:        http.StatusServiceUnavailable,
		Message:     "Base de données indisponible",
		Details:     "Impossible de se connecter à la base de données",
		InternalErr: err,
	}
}
