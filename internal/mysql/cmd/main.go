package main

import (
	"fmt"
	"os"

	"modules-minibackup/internal/mysql/pkg/commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mini-backup",
	Short: "Un outil CLI pour gérer les backups MySQL",
	Long:  `Un module CLI basé sur Cobra.`,
}

func main() {
	rootCmd.AddCommand(commands.BackupCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
