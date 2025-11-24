package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// Done LinkRepository est une interface qui définit les méthodes d'accès aux données
// pour les opérations CRUD sur les liens.
// L'implémenter avec les méthodes nécessaires
type LinkRepository interface {
	// CreateLink insère un nouveau lien dans la base de données.
	CreateLink(link *models.Link) error
	// DeleteLink supprime un lien de la base de données en utilisant son ID.
	DeleteLink(linkID uint) error
	// GetLinkByShortCode récupère un lien de la base de données en utilisant son shortCode.
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	// GetAllLinks récupère tous les liens de la base de données.
	GetAllLinks() ([]models.Link, error)
	// CountClicksByLinkID compte le nombre total de clics pour un ID de lien donné.
	CountClicksByLinkID(linkID uint) (int, error)
}

// Done :  GormLinkRepository est l'implémentation de LinkRepository utilisant GORM.
type GormLinkRepository struct {
	db *gorm.DB
}

// NewLinkRepository crée et retourne une nouvelle instance de GormLinkRepository.
// Cette fonction retourne *GormLinkRepository, qui implémente l'interface LinkRepository.
func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	// Done
	return &GormLinkRepository{db: db}
}

// CreateLink insère un nouveau lien dans la base de données.
func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	// Done 1: Utiliser GORM pour créer un nouvel enregistrement (link) dans la table des liens.
	result := r.db.Create(link)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteLink supprime un lien de la base de données en utilisant son ID.
func (r *GormLinkRepository) DeleteLink(linkID uint) error {
	// Utiliser GORM pour supprimer le lien avec l'ID donné.
	result := r.db.Delete(&models.Link{}, linkID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetLinkByShortCode récupère un lien de la base de données en utilisant son shortCode.
// Il renvoie gorm.ErrRecordNotFound si aucun lien n'est trouvé avec ce shortCode.
func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	// Done 2: Utiliser GORM pour trouver un lien par son ShortCode.
	// La méthode First de GORM recherche le premier enregistrement correspondant et le mappe à 'link'.
	result := r.db.Where("shortcode = ?", shortCode).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

// GetAllLinks récupère tous les liens de la base de données.
// Cette méthode est utilisée par le moniteur d'URLs.
func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	// Done 3: Utiliser GORM pour récupérer tous les liens.
	result := r.db.Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}

// CountClicksByLinkID compte le nombre total de clics pour un ID de lien donné.
func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64 // GORM retourne un int64 pour les comptes
	// Done 4: Utiliser GORM pour compter les enregistrements dans la table 'clicks'
	// où 'LinkID' correspond à l'ID du lien donné.
	result := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}
