package services

import "errors"

// Erreurs spécifiques du service
var (
	// ErrURLAlreadyExists est retourné quand une URL longue existe déjà dans la base
	ErrURLAlreadyExists = errors.New("URL existante")

	// ErrShortCodeCollision est retourné quand tous les codes courts générés sont déjà utilisés
	ErrShortCodeCollision = errors.New("failed to generate unique short code after maximum retries")
)
