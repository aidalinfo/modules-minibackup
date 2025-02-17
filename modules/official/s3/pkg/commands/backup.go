package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/s3/pkg/utils"

	"github.com/spf13/cobra"
)

func BackupCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "backup [name] [args]",
		Short: "Exécute un backup S3 avec des paramètres JSON",
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

			result, err := utils.BackupRemoteS3(name, backupArgs, loggerModule)
			if err != nil {
				loggerModule.Error(fmt.Sprintf("Erreur lors du backup : %v", err))
				log.Fatalf("❌ Erreur lors du backup : %v", err)
			}
			loggerModule.Info("Backup S3 exécuté avec succès.")
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
