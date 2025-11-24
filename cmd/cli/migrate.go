package cli

import (
	"fmt"
	"log"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"

	// "github.com/glebarez/sqlite" WINDOWS
	"gorm.io/driver/sqlite" // MAC
	"gorm.io/gorm"
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// DONE : Charger la configuration chargée globalement via cmd.Cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Configuration non chargée")
		}

		// DONE : Initialiser la connexion à la BDD
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		// DONE : Assurez-vous que la connexion est fermée après la migration grâce à defer
		defer sqlDB.Close()

		// DONE : Exécuter les migrations automatiques de GORM.
		// Utilisez db.AutoMigrate() et passez-lui les pointeurs vers tous vos modèles.
		log.Println("Exécution des migrations de la base de données...")
		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("FATAL: Échec des migrations: %v", err)
		}

		// Pas touche au log
		fmt.Println("Migrations de la base de données exécutées avec succès.\n")
	},
}

func init() {
	// DONE : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(MigrateCmd)
}
