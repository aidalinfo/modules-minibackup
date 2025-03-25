package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RestoreFolder restaure les dossiers définis dans config.Fs.Paths.
func RestoreFolder(restorePath string, config BackupArgs) error {
	fmt.Printf("Début de la restauration depuis : %s\n", restorePath)
	fmt.Printf("Chemins à restaurer : %v\n", config.Fs.Paths)

	for _, destPath := range config.Fs.Paths {
		os.Mkdir(destPath, 0755)

		baseName := filepath.Base(destPath)
		pattern := filepath.Join(restorePath, baseName+"_*")
		fmt.Printf("Recherche de sauvegarde pour %s avec le pattern %s\n", destPath, pattern)

		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("erreur lors de la recherche avec Glob pour %s : %v", pattern, err)
		}
		if len(matches) == 0 {
			fmt.Printf("Aucune sauvegarde trouvée pour %s\n", destPath)
			continue
		}
		sourcePath := matches[0]

		fmt.Printf("Restauration de %s vers %s\n", sourcePath, destPath)
		err = os.MkdirAll(destPath, 0755)
		if err != nil {
			return fmt.Errorf("erreur création dossier destination %s : %v", destPath, err)
		}

		err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsPermission(err) || strings.Contains(err.Error(), "permission denied") {
					fmt.Printf("Permission refusée, skip : %s\n", path)
					return nil
				}
				return err
			}

			rel, err := filepath.Rel(sourcePath, path)
			if err != nil {
				return fmt.Errorf("erreur calcul chemin relatif : %v", err)
			}
			targetPath := filepath.Join(destPath, rel)

			if info.IsDir() {
				if mkErr := os.MkdirAll(targetPath, info.Mode()); mkErr != nil {
					if os.IsPermission(mkErr) {
						fmt.Printf("Permission refusée sur dossier, skip : %s\n", targetPath)
						return nil
					}
					return mkErr
				}
				return nil
			}

			if cpErr := copyFile(path, targetPath); cpErr != nil {
				if os.IsPermission(cpErr) || strings.Contains(cpErr.Error(), "permission denied") {
					fmt.Printf("Permission refusée sur fichier, skip : %s\n", path)
					return nil
				}
				return cpErr
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("erreur lors de la restauration de %s : %v", destPath, err)
		}
		fmt.Printf("Restauration terminée pour %s\n", destPath)
	}

	return nil
}

// copyFile copie un fichier d'une source vers une destination.
func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Permission refusée ouverture : %s, skip\n", srcPath)
			return nil
		}
		return fmt.Errorf("erreur ouverture source %s : %v", srcPath, err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Permission refusée création : %s, skip\n", destPath)
			return nil
		}
		return fmt.Errorf("erreur création destination %s : %v", destPath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Permission refusée copy : %s, skip\n", srcPath)
			return nil
		}
		return fmt.Errorf("erreur copie de %s vers %s : %v", srcPath, destPath, err)
	}

	fmt.Printf("Fichier %s copié avec succès vers %s\n", srcPath, destPath)
	return nil
}
