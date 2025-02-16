package main

import (
	"modules-minibackup/internal/mysql/pkg/commands"

	"github.com/spf13/cobra"
)

func main() {
	// Crée le rootCmd et y ajoute les commandes backup et restore.
	rootCmd := &cobra.Command{
		Use:   "module-mysql-mb",
		Short: "Module MySQL pour MiniBackup",
		Long: `
Ce module permet de gérer les backups et la restauration.
`,
	}
	rootCmd.AddCommand(commands.BackupCmd())
	rootCmd.AddCommand(commands.RestoreCmd())

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
