package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Int64Ptr(i int64) *int64 {
	return &i
}

type S3Manager struct {
	Client *s3.Client
	Bucket string
}

type BackupDetails struct {
	Key          string
	Size         int64
	LastModified time.Time
}

func AwsCredentialFileCreateFunc(accessKey, secretKey string, header string) error {
	// Définir le chemin du fichier credentials
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du répertoire personnel : %v", err)
	}
	awsCredentialsPath := filepath.Join(homeDir, ".aws", "credentials")

	// Créer le dossier ~/.aws s'il n'existe pas
	err = os.MkdirAll(filepath.Dir(awsCredentialsPath), 0700)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du dossier .aws : %v", err)
	}

	// Lire le contenu existant du fichier credentials s'il existe
	var existingContent string
	if _, err := os.Stat(awsCredentialsPath); err == nil {
		data, err := os.ReadFile(awsCredentialsPath)
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture du fichier credentials : %v", err)
		}
		existingContent = string(data)
	}
	var sectionHeader string
	if header == "" {
		sectionHeader = "[default-s3modules-minibackup]"
	} else {
		sectionHeader = "[" + header + "]"
	}
	if existingContent != "" && containsSection(existingContent, sectionHeader) {
		fmt.Println("La section " + sectionHeader + " existe déjà dans le fichier credentials.")
		return nil
	}

	newSection := fmt.Sprintf(`%s
aws_access_key_id = %s
aws_secret_access_key = %s
`, sectionHeader, accessKey, secretKey)

	newContent := existingContent + "\n" + newSection

	// Écrire le contenu mis à jour dans le fichier credentials
	err = os.WriteFile(awsCredentialsPath, []byte(newContent), 0600)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier credentials : %v", err)
	}

	fmt.Println("Fichier credentials AWS créé avec succès.")
	return nil
}

// containsSection vérifie si une section existe déjà dans le contenu du fichier
func containsSection(content, sectionHeader string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == sectionHeader {
			return true
		}
	}
	return false
}

// NewS3Manager initialise le gestionnaire S3 en utilisant la configuration AWS par défaut
func NewS3Manager(bucket, region, endpoint string, awsprofile string, pathStyle bool) (*S3Manager, error) {
	// Charger la configuration par défaut depuis les fichiers AWS (credentials et config)
	var profileName string
	if awsprofile == "" {
		profileName = "default-s3modules-minibackup"
	} else {
		profileName = awsprofile
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region), // Région par défaut
		config.WithSharedConfigProfile(profileName),
	)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du chargement de la configuration AWS : %v", err)
	}
	// Initialiser le client S3 avec le point de terminaison Scaleway
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = pathStyle      // Mode de chemin d'accès (obligatoire pour Scaleway)
		o.BaseEndpoint = &endpoint // Point de terminaison personnalisé
	})

	return &S3Manager{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (m *S3Manager) ListBuckets() ([]string, error) {
	// Récupération de la liste des buckets
	result, err := m.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("Erreur lors de la récupération des buckets : %v", err)
		return nil, fmt.Errorf("erreur lors de la récupération des buckets : %v", err)
	}

	// Extraction des noms des buckets
	buckets := []string{}
	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}

	fmt.Printf("Liste des buckets S3 récupérée avec succès : %v", buckets)
	return buckets, nil
}

// copyBackupToLocal copie tout le contenu d'un bucket S3 vers un répertoire local
func (m *S3Manager) CopyBackupToLocal(destination string) error {
	// Lister tous les objets dans le bucket
	listInput := &s3.ListObjectsV2Input{
		Bucket: &m.Bucket,
	}

	// Obtenir la liste des objets
	result, err := m.Client.ListObjectsV2(context.TODO(), listInput)
	if err != nil {
		return fmt.Errorf("erreur lors de la liste des objets dans le bucket %s : %v", m.Bucket, err)
	}

	// Parcourir chaque objet dans le bucket
	for _, object := range result.Contents {
		// Calculer le chemin local correspondant
		localPath := filepath.Join(destination, *object.Key)

		// Créer les répertoires nécessaires pour le fichier
		err := os.MkdirAll(filepath.Dir(localPath), 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du répertoire pour %s : %v", localPath, err)
		}

		// Télécharger l'objet
		getInput := &s3.GetObjectInput{
			Bucket: &m.Bucket,
			Key:    object.Key,
		}
		objectOutput, err := m.Client.GetObject(context.TODO(), getInput)
		if err != nil {
			return fmt.Errorf("erreur lors du téléchargement de l'objet %s : %v", *object.Key, err)
		}
		defer objectOutput.Body.Close()

		// Créer le fichier local
		localFile, err := os.Create(localPath)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du fichier local %s : %v", localPath, err)
		}
		defer localFile.Close()

		// Copier le contenu de l'objet dans le fichier local
		_, err = io.Copy(localFile, objectOutput.Body)
		if err != nil {
			return fmt.Errorf("erreur lors de la copie du contenu de %s vers %s : %v", *object.Key, localPath, err)
		}

		fmt.Printf("Backup copié avec succès depuis %s vers %s\n", *object.Key, localPath)
	}

	return nil
}
