package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"modules-minibackup/internal/mongo/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

func RestoreCmd() *cobra.Command {
	var argsJSON string

	cmd := &cobra.Command{
		Use:   "restore [name] [backupPath] [args]",
		Short: "Exécute une restauration Mongo avec des paramètres JSON",
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

			err := utils.RestoreMongoDB(backupPath, restoreArgs, loggerModule)
			if err != nil {
				loggerModule.Error(fmt.Sprintf("Erreur lors de la restauration : %v", err))
				log.Fatalf("❌ Erreur lors de la restauration : %v", err)
				os.Exit(1)
			}
			loggerModule.Info("Restauration Mongo exécutée avec succès.")
			fmt.Println(true)
		},
	}

	return cmd
}
