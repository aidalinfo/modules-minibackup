package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/s3/pkg/utils"

	"github.com/spf13/cobra"
)

func RestoreCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "restore [name] [backupPath] [args]",
		Short: "Exécute une restauration S3 avec des paramètres JSON",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			backupPath := args[1]
			argsJSON = args[2]

			loggerModule := utils.NewModuleLogger()
			loggerModule.Info(fmt.Sprintf("Démarrage de la restauration pour %s", name))

			var restoreArgs utils.BackupArgs
			if err := json.Unmarshal([]byte(argsJSON), &restoreArgs); err != nil {
				loggerModule.Error(fmt.Sprintf("Erreur de parsing JSON: %v", err))
				log.Fatalf("❌ Erreur de parsing JSON: %v", err)
			}
			loggerModule.Info("Arguments de restauration parsés avec succès.")

			err := utils.RestoreS3(backupPath, restoreArgs, name, loggerModule)
			if err != nil {
				loggerModule.Error(fmt.Sprintf("Erreur lors de la restauration : %v", err))
				log.Fatalf("❌ Erreur lors de la restauration : %v", err)
				loggerModule.SetResult(false)
				jsonOutput, err := loggerModule.JSON()
				if err != nil {
					log.Fatalf("Erreur lors de la sérialisation du logger en JSON: %v", err)
				}
				fmt.Println(jsonOutput)
				return
			}
			loggerModule.Info("Restauration S3 exécutée avec succès.")
			loggerModule.SetResult(true)
			jsonOutput, err := loggerModule.JSON()
			if err != nil {
				log.Fatalf("Erreur lors de la sérialisation du logger en JSON: %v", err)
			}
			fmt.Println(jsonOutput)
		},
	}

	return cmd
}
