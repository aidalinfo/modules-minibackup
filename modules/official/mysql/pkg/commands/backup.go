package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/mysql/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

func BackupCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "backup [name] [args]",
		Short: "Exécute un backup MySQL avec des paramètres JSON",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			argsJSON = args[1]

			loggerModule := utils.NewModuleLogger()
			loggerModule.Info(fmt.Sprintf("Démarrage du backup pour %s", name))

			var backupArgs utils.BackupArgs
			err := json.Unmarshal([]byte(argsJSON), &backupArgs)
			if err != nil {
				loggerModule.Error(fmt.Sprintf("Erreur de parsing JSON: %v", err))
			}
			loggerModule.Info("Arguments de backup parsés avec succès.")

			result, err := utils.BackupMySQL(name, backupArgs, loggerModule)
			if err != nil {
				loggerModule.Error(fmt.Sprintf("Erreur lors du backup : %v", err))
				log.Fatalf("❌ Erreur lors du backup : %v", err)
				os.Exit(1)
			}
			loggerModule.Info("Backup MySQL exécuté avec succès.")
			fmt.Println(result)
		},
	}

	return cmd
}
