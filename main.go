package main

import (
	"github.com/axellelanca/urlshortener/cmd"
	_ "github.com/axellelanca/urlshortener/cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "github.com/axellelanca/urlshortener/cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
)

func main() {
	// DONE
	// Exécute la commande racine Cobra
	// Cela déclenchera l'analyse des arguments CLI et l'exécution de la commande appropriée
	cmd.Execute()
}
