package cli

import (
	"errors"
	"fmt"
	"log"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : variable shortCodeFlag qui stockera la valeur du flag --code
var shortCodeFlag string

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		// DONE : Valider que le flag --code a été fourni.
		if shortCodeFlag == "" {
			log.Fatal("FATAL: Le flag --code est requis")
		}

		// DONE : Charger la configuration chargée globalement via cmd.Cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatal("FATAL: Configuration non chargée")
		}

		// DONE : Initialiser la connexion à la BDD.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// DONE S'assurer que la connexion est fermée à la fin de l'exécution de la commande grâce à defer
		defer sqlDB.Close()

		// DONE : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// DONE : Appeler GetLinkStats pour récupérer le lien et ses statistiques.
		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("ERREUR: Aucun lien trouvé avec le code '%s'", shortCodeFlag)
			}
			log.Fatalf("FATAL: Erreur lors de la récupération des statistiques: %v", err)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.Shortcode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n\n", totalClicks)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// DONE : Définir le flag --code pour la commande stats.
	StatsCmd.Flags().StringVarP(&shortCodeFlag, "code", "c", "", "Code court de l'URL")

	// DONE Marquer le flag comme requis
	StatsCmd.MarkFlagRequired("code")

	// DONE : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(StatsCmd)
}
