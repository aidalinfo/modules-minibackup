package main

import (
	"fmt"
	"os"

	"modules-minibackup/internal/mysql/pkg/commands"

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
}
