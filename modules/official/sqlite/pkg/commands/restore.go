package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/sqlite/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

func RestoreCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "restore [name] [backupPath] [args]",
		Short: "Exécute une restauration sqlite avec des paramètres JSON",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			backupPath := args[1]
			argsJSON = args[2]

			fmt.Println(name)
			var restoreArgs utils.BackupArgs
			if err := json.Unmarshal([]byte(argsJSON), &restoreArgs); err != nil {
				log.Fatalf("❌ Erreur de parsing JSON: %v", err)
			}

			err := utils.RestoreSqlite(name, restoreArgs, backupPath)
			if err != nil {
				log.Fatalf("❌ Erreur lors de la restauration : %v", err)
				os.Exit(1)
			}
			fmt.Println(true)
		},
	}

	return cmd
}
