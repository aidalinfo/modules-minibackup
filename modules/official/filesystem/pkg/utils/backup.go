package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CopyFolder copie un ou plusieurs dossiers et leur contenu vers un dossier de destination.
func CopyFolder(name string, config BackupArgs) (string, error) {
    paths := config.Fs.Paths
    date := time.Now().Format("20060102_150405")
    
    absParentDir, err := filepath.Abs(fmt.Sprintf("%s/%s_fs_backup_%s", config.Path, name, date))
    if err != nil {
        return "", fmt.Errorf("erreur lors de la conversion du chemin absolu : %v", err)
    }
    
    err = os.MkdirAll(absParentDir, 0755)
    if err != nil {
        return "", fmt.Errorf("erreur lors de la création du dossier de destination %s : %v", absParentDir, err)
    }
    fmt.Printf("Dossier de destination créé avec succès : %s\n", absParentDir)
    fmt.Printf("Liste des chemins à sauvegarder : %v\n", paths)

    for _, srcPath := range paths {
        baseName := filepath.Base(srcPath)
        newFolderName := fmt.Sprintf("%s_%s", baseName, date)
        newDestination := filepath.Join(absParentDir, newFolderName)

        info, err := os.Stat(srcPath)
        if err != nil {
            fmt.Printf("Erreur sur %s : %v\n", srcPath, err)
            continue
        }

        if info.IsDir() {
            fmt.Printf("Copie du dossier %s vers %s\n", srcPath, newDestination)
            
            err := CopyDir(newDestination, srcPath)
            if err != nil {
                fmt.Printf("Erreur lors de la copie de %s : %v\n", srcPath, err)
                continue
            }
            fmt.Printf("Dossier %s copié avec succès\n", srcPath)
        } else {
            fmt.Printf("Copie du fichier %s vers %s\n", srcPath, newDestination)
            err := CopyFile(srcPath, newDestination)
            if err != nil {
                fmt.Printf("Erreur lors de la copie du fichier %s : %v\n", srcPath, err)
                continue
            }
            fmt.Printf("Fichier %s copié avec succès\n", srcPath)
        }
    }


    return absParentDir, nil
}

func CopyFile(srcFile, destFile string) error {
	// Assure la création des dossiers parents
	err := os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return fmt.Errorf("erreur création du dossier parent %s : %v", filepath.Dir(destFile), err)
	}

	src, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("erreur ouverture source %s : %v", srcFile, err)
	}
	defer src.Close()

	dst, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("erreur création destination %s : %v", destFile, err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("erreur copie %s vers %s : %v", srcFile, destFile, err)
	}
	return nil
}
