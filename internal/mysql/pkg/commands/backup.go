package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/mysql/pkg/utils"

	"github.com/spf13/cobra"
)

// BackupCmd crée la commande CLI pour le backup MySQL
func BackupCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "backup [name] [args]",
		Short: "Exécute un backup MySQL avec des paramètres JSON",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			argsJSON = args[1]

			var backupArgs utils.BackupArgs
			err := json.Unmarshal([]byte(argsJSON), &backupArgs)
			if err != nil {
				log.Fatalf("❌ Erreur de parsing JSON: %v", err)
			}

			result, err := utils.BackupMySQL(name, backupArgs)
			if err != nil {
				log.Fatalf("❌ Erreur lors du backup : %v", err)
			}

			fmt.Println("✅ Backup réussi:", result)
		},
	}

	return cmd
}
