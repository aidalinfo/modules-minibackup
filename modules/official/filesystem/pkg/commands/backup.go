package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/filesystem/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

func BackupCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "backup [name] [args]",
		Short: "Exécute un backup Filesystem avec des paramètres JSON",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			argsJSON = args[1]

			var backupArgs utils.BackupArgs
			err := json.Unmarshal([]byte(argsJSON), &backupArgs)
			if err != nil {
				log.Fatalf("❌ Erreur de parsing JSON: %v", err)
			}

			result, err := utils.CopyFolder(name, backupArgs)
			if err != nil {
				log.Fatalf("❌ Erreur lors du backup : %v", err)
				os.Exit(1)
			}
			log.Printf("Backup Filesystem exécuté avec succès.")
			fmt.Println(result)
		},
	}

	return cmd
}
