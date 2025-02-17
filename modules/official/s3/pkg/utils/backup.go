package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func BackupRemoteS3(name string, config BackupArgs, logger *ModuleLogger) ([]string, error) {
	logger.Debug(fmt.Sprintf("Information de connexion S3 : %v", config.S3))

	// Création du fichier credentials AWS
	err := AwsCredentialFileCreateFunc(config.S3.ACCESS_KEY, config.S3.SECRET_KEY, name)
	if err != nil {
		return nil, err
	}

	// Formatage du dossier parent avec timestamp
	date := time.Now().Format("20060102_150405")
	parentDir := fmt.Sprintf("%s/%s_s3_backup_%s", config.Path, name, date)

	// Liste des buckets à sauvegarder
	var bucketsToBackup []string

	// Si `All` est activé, lister tous les buckets disponibles
	if config.S3.All {
		s3client, err := NewS3Manager("", config.S3.Region, config.S3.Endpoint, name, config.S3.PathStyle)
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de l'initialisation du gestionnaire S3 : %v", err))
			return nil, err
		}

		buckets, err := s3client.ListBuckets()
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de la récupération de la liste des buckets S3 : %v", err))
			return nil, err
		}

		bucketsToBackup = buckets
	} else {
		bucketsToBackup = config.S3.Bucket
	}

	// Liste des chemins de backup
	allBucketPath := []string{}

	// Sauvegarde chaque bucket
	for _, bucket := range bucketsToBackup {
		logger.Debug(fmt.Sprintf("Backup du bucket S3 : %s", bucket))

		s3client, err := NewS3Manager(bucket, config.S3.Region, config.S3.Endpoint, name, config.S3.PathStyle)
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de l'initialisation du gestionnaire S3 pour %s : %v", bucket, err))
			continue
		}

		// Destination du backup local
		destinationPath := filepath.Join(parentDir, fmt.Sprintf("%s-%s-%s", name, bucket, date))

		// Copier le backup depuis S3
		err = s3client.CopyBackupToLocal(destinationPath)
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de la copie du backup depuis S3 pour %s : %v", bucket, err))
			continue
		}

		allBucketPath = append(allBucketPath, destinationPath)
	}

	logger.Info(fmt.Sprintf("Backup copié avec succès depuis les buckets S3 : %v", allBucketPath))
	return []string{parentDir}, nil
}

// UploadEmptyFolder crée un "dossier" vide dans S3
func (m *S3Manager) UploadEmptyFolder(folderPath string) error {
	// Ajouter "/" à la fin pour indiquer un dossier
	if !strings.HasSuffix(folderPath, "/") {
			folderPath += "/"
	}

	input := &s3.PutObjectInput{
			Bucket:      &m.Bucket,
			Key:         &folderPath,
			ContentType: aws.String("application/x-directory"), // MIME type indiquant un dossier
	}

	_, err := m.Client.PutObject(context.TODO(), input)
	if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier %s : %v", folderPath, err)
	}

	fmt.Sprintf("Dossier %s créé avec succès dans S3", folderPath)
	return nil
}

// UploadFileToS3 téléverse un fichier local vers un chemin S3
func (m *S3Manager) Upload(localPath, s3Path string, useGlacier bool) error {
	// Ouvrir le fichier local
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier %s : %v", localPath, err)
	}
	defer file.Close()

	// Obtenir la taille du fichier
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des informations du fichier %s : %v", localPath, err)
	}
	// Déterminer la classe de stockage (STANDARD ou GLACIER)
	var storageClass types.StorageClass = types.StorageClassStandard // Valeur par défaut
	if useGlacier {
		storageClass = types.StorageClassGlacier
	}

	// Préparer la requête de téléversement
	input := &s3.PutObjectInput{
		Bucket:        &m.Bucket,
		Key:           &s3Path,
		Body:          file,
		ContentLength: Int64Ptr(stat.Size()), // Convertir en *int64
		ContentType:   aws.String("application/octet-stream"),
		StorageClass:  storageClass,
	}

	// Téléverser le fichier
	_, err = m.Client.PutObject(context.TODO(), input)
	if err != nil {
		// getLogger().Error(fmt.Sprintf("Erreur lors de la téléversement du fichier %s vers %s : %v", localPath, s3Path, err))
		return fmt.Errorf("erreur lors de l'upload vers S3 (local: %s, s3: %s) : %v", localPath, s3Path, err)
	}
	fmt.Sprintf("Fichier %s téléversé avec succès vers %s", localPath, s3Path)
	// getLogger().Info(fmt.Sprintf("Fichier %s téléversé avec succès vers %s", localPath, s3Path))
	return nil
}

func (m *S3Manager) DoesBucketExist(bucketName string) (bool, error) {
	_, err := m.Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		var notFoundErr *types.NotFound
		if errors.As(err, &notFoundErr) {
			return false, nil // Le bucket n'existe pas
		}
		return false, err // Autre erreur
	}
	return true, nil // Le bucket existe
}

// CreateBucket crée un bucket S3 s'il n'existe pas déjà
func (m *S3Manager) CreateBucket(bucketName string) error {
	// Préparer l'entrée pour créer un bucket
	input := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	// Création du bucket
	_, err := m.Client.CreateBucket(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du bucket %s : %v", bucketName, err)
	}

	fmt.Sprintln("Bucket créé avec succès")
	return nil
}
