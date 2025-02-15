package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/mysql/pkg/utils"

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
			loggerModule.AddInfo(fmt.Sprintf("Démarrage du backup pour %s", name))

			var backupArgs utils.BackupArgs
			err := json.Unmarshal([]byte(argsJSON), &backupArgs)
			if err != nil {
				loggerModule.AddError(fmt.Sprintf("Erreur de parsing JSON: %v", err))
				log.Fatalf("❌ Erreur de parsing JSON: %v", err)
			}
			loggerModule.AddInfo("Arguments de backup parsés avec succès.")

			result, err := utils.BackupMySQL(name, backupArgs)
			if err != nil {
				loggerModule.AddError(fmt.Sprintf("Erreur lors du backup : %v", err))
				log.Fatalf("❌ Erreur lors du backup : %v", err)
			}
			loggerModule.AddInfo("Backup MySQL exécuté avec succès.")

			loggerModule.SetResult(result)

			jsonOutput, err := loggerModule.JSON()
			if err != nil {
				log.Fatalf("Erreur lors de la sérialisation du logger en JSON: %v", err)
			}
			fmt.Println(jsonOutput)
		},
	}

	return cmd
}
