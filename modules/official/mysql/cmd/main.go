package main

import (
	"fmt"
	"modules-minibackup/internal/mysql/pkg/commands"
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "module-mysql-mb",
	Short: "Module MySQL pour MiniBackup",
	Long: `
Ce module permet de gérer les backups et la restauration.
`,
}

func main() {

	rootCmd.AddCommand(commands.BackupCmd())
	rootCmd.AddCommand(commands.RestoreCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // Crée le dossier "docs" pour stocker la doc
	// if err := os.MkdirAll("docs", 0755); err != nil {
	// 	log.Fatalf("Erreur lors de la création du dossier docs: %v", err)
	// }

	// // Génère la documentation Markdown pour toutes les commandes
	// if err := doc.GenMarkdownTree(rootCmd, "docs"); err != nil {
	// 	log.Fatalf("Erreur lors de la génération de la documentation: %v", err)
	// }

	// log.Println("Documentation générée dans le dossier docs")
}
