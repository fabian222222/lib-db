package database

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupMetadata struct {
	DatabaseName   string    `json:"database_name"`
	OriginalOwner  string    `json:"original_owner"`
	BackupDate     time.Time `json:"backup_date"`
	LibDBVersion   string    `json:"lib_db_version"`
}

func CreateBackup(dbName, backupFile, owner string) error {
	dbPath := filepath.Join("../../databases", dbName)
	
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("la base de données '%s' n'existe pas", dbName)
	}

	zipFile, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	metadata := BackupMetadata{
		DatabaseName:  dbName,
		OriginalOwner: owner,
		BackupDate:    time.Now(),
		LibDBVersion:  "1.0",
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataFile, err := zipWriter.Create(".backup_metadata.json")
	if err != nil {
		return err
	}
	_, err = metadataFile.Write(metadataBytes)
	if err != nil {
		return err
	}

	err = filepath.Walk(dbPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dbPath, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipFileEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFileEntry, file)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func CanUserRestoreBackup(backupFile, currentUser string) (bool, string, error) {
	zipReader, err := zip.OpenReader(backupFile)
	if err != nil {
		return false, "", fmt.Errorf("impossible d'ouvrir le fichier de sauvegarde: %v", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == ".backup_metadata.json" {
			fileReader, err := file.Open()
			if err != nil {
				return false, "", fmt.Errorf("impossible de lire les métadonnées: %v", err)
			}
			defer fileReader.Close()

			var metadata BackupMetadata
			decoder := json.NewDecoder(fileReader)
			if err := decoder.Decode(&metadata); err != nil {
				return false, "", fmt.Errorf("métadonnées invalides: %v", err)
			}

			if currentUser == "admin" || metadata.OriginalOwner == currentUser {
				return true, metadata.OriginalOwner, nil
			}
			return false, metadata.OriginalOwner, nil
		}
	}

	return false, "", fmt.Errorf("métadonnées de sauvegarde non trouvées")
}

func RestoreBackup(backupFile, newDbName string) error {
	zipReader, err := zip.OpenReader(backupFile)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir le fichier de sauvegarde: %v", err)
	}
	defer zipReader.Close()

	dbPath := filepath.Join("../../databases", newDbName)
	
	if _, err := os.Stat(dbPath); err == nil {
		return fmt.Errorf("la base de données '%s' existe déjà", newDbName)
	}

	err = os.MkdirAll(dbPath, 0755)
	if err != nil {
		return fmt.Errorf("impossible de créer le dossier de la base: %v", err)
	}

	for _, file := range zipReader.File {
		if file.Name == ".backup_metadata.json" {
			continue
		}

		path := filepath.Join(dbPath, file.Name)
		dir := filepath.Dir(path)
		
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("impossible de créer le dossier: %v", err)
		}

		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("impossible d'ouvrir le fichier dans l'archive: %v", err)
		}
		defer fileReader.Close()

		outFile, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("impossible de créer le fichier: %v", err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, fileReader)
		if err != nil {
			return fmt.Errorf("impossible de copier le fichier: %v", err)
		}
	}

	return nil
}

func GetBackupInfo(backupFile string) (*BackupMetadata, error) {
	zipReader, err := zip.OpenReader(backupFile)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier de sauvegarde: %v", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == ".backup_metadata.json" {
			fileReader, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("impossible de lire les métadonnées: %v", err)
			}
			defer fileReader.Close()

			var metadata BackupMetadata
			decoder := json.NewDecoder(fileReader)
			if err := decoder.Decode(&metadata); err != nil {
				return nil, fmt.Errorf("métadonnées invalides: %v", err)
			}

			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("métadonnées de sauvegarde non trouvées")
} 