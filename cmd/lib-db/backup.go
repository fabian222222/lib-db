package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"github.com/fabian222222/lib-db/pkg/database"
)

type BackupMetadata struct {
	DatabaseName   string    `json:"database_name"`
	OriginalOwner  string    `json:"original_owner"`
	BackupDate     time.Time `json:"backup_date"`
	LibDBVersion   string    `json:"lib_db_version"`
}

func handleBackup(args []string) {
	if len(args) >= 1 && args[0] == "info" {
		handleBackupInfo(args[1:])
		return
	}

	ok, session, err := database.IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour faire une sauvegarde.")
		return
	}

	if len(args) < 1 {
		fmt.Println("Usage : backup <database_name> [backup_file]")
		fmt.Println("       : backup info <backup_file>")
		return
	}

	dbName := args[0]
	
	if !database.UserHasAccess(session.Username, dbName) {
		fmt.Printf("❌ Vous n'avez pas les permissions pour sauvegarder la base '%s'.\n", dbName)
		fmt.Println("💡 Seuls les propriétaires peuvent sauvegarder leurs bases de données.")
		return
	}
	
	backupDir := "../../backup"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err = os.MkdirAll(backupDir, 0755)
		if err != nil {
			fmt.Println("Erreur lors de la création du dossier backup :", err)
			return
		}
	}
	
	backupFile := fmt.Sprintf("backup_%s_%s.zip", dbName, time.Now().Format("20060102_150405"))
	
	if len(args) > 1 {
		backupFile = args[1]
		if filepath.Ext(backupFile) != ".zip" {
			backupFile += ".zip"
		}
	}
	
	backupFile = filepath.Join(backupDir, backupFile)

	err = createBackup(dbName, backupFile, session.Username)
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde :", err)
		return
	}

	fmt.Printf("✅ Sauvegarde créée : %s\n", backupFile)
	fmt.Printf("🔒 Seul %s pourra restaurer cette sauvegarde.\n", session.Username)
}

func handleRestore(args []string) {
	ok, session, err := database.IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour restaurer une base de données.")
		return
	}

	if len(args) < 2 {
		fmt.Println("Usage : restore <backup_file> <new_database_name>")
		return
	}

	backupFile := args[0]
	newDbName := args[1]
	
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		backupInDir := filepath.Join("../../backup", backupFile)
		if _, err := os.Stat(backupInDir); err == nil {
			backupFile = backupInDir
		}
	}
	
	canRestore, originalOwner, err := canUserRestoreBackup(backupFile, session.Username)
	if err != nil {
		fmt.Println("Erreur lors de la vérification du backup :", err)
		return
	}
	
	if !canRestore {
		fmt.Printf("❌ Vous ne pouvez pas restaurer cette sauvegarde.\n")
		fmt.Printf("🔒 Cette sauvegarde appartient à : %s\n", originalOwner)
		fmt.Println("💡 Seul le propriétaire original peut restaurer ses sauvegardes.")
		return
	}
	
	fmt.Printf("🔄 Restauration de '%s' vers la nouvelle base '%s'\n", backupFile, newDbName)
	fmt.Printf("✅ Vérification des permissions : OK (propriétaire : %s)\n", session.Username)

	err = restoreBackup(backupFile, newDbName)
	if err != nil {
		fmt.Println("Erreur lors de la restauration :", err)
		return
	}

	// Accorder l'accès à la base restaurée (puisque c'est le propriétaire original)
	err = database.GrantDatabaseAccess(session.Username, newDbName)
	if err != nil {
		fmt.Printf("⚠️ Base restaurée mais erreur d'attribution des droits : %v\n", err)
	} else {
		fmt.Printf("✅ Base de données '%s' restaurée avec succès !\n", newDbName)
		fmt.Printf("🔑 Vous récupérez vos droits de propriétaire.\n")
	}
}

func createBackup(dbName, backupFile, owner string) error {
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
	
	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
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

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dbPath, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func canUserRestoreBackup(backupFile, currentUser string) (bool, string, error) {
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return false, "", fmt.Errorf("le fichier de sauvegarde '%s' n'existe pas", backupFile)
	}

	zipReader, err := zip.OpenReader(backupFile)
	if err != nil {
		return false, "", err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == ".backup_metadata.json" {
			fileReader, err := file.Open()
			if err != nil {
				return false, "", err
			}
			defer fileReader.Close()

			var metadata BackupMetadata
			decoder := json.NewDecoder(fileReader)
			err = decoder.Decode(&metadata)
			if err != nil {
				return false, "", err
			}

			return metadata.OriginalOwner == currentUser, metadata.OriginalOwner, nil
		}
	}

	return false, "unknown", fmt.Errorf("cette sauvegarde ne contient pas de métadonnées de propriétaire")
}

func handleBackupInfo(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : backup info <backup_file>")
		return
	}

	backupFile := args[0]
	
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		backupInDir := filepath.Join("../../backup", backupFile)
		if _, err := os.Stat(backupInDir); err == nil {
			backupFile = backupInDir
		} else {
			fmt.Printf("❌ Le fichier '%s' n'existe pas.\n", args[0])
			return
		}
	}

	zipReader, err := zip.OpenReader(backupFile)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == ".backup_metadata.json" {
			fileReader, err := file.Open()
			if err != nil {
				fmt.Println("Erreur lors de la lecture des métadonnées :", err)
				return
			}
			defer fileReader.Close()

			var metadata BackupMetadata
			decoder := json.NewDecoder(fileReader)
			err = decoder.Decode(&metadata)
			if err != nil {
				fmt.Println("Erreur lors du décodage des métadonnées :", err)
				return
			}

			fmt.Printf("\n📦 INFORMATIONS DE SAUVEGARDE\n")
			fmt.Println("═════════════════════════════")
			fmt.Printf("📁 Base de données : %s\n", metadata.DatabaseName)
			fmt.Printf("👤 Propriétaire : %s\n", metadata.OriginalOwner)
			fmt.Printf("📅 Date de sauvegarde : %s\n", metadata.BackupDate.Format("02/01/2006 15:04:05"))
			fmt.Printf("🔧 Version Lib-DB : %s\n", metadata.LibDBVersion)
			
			if info, err := os.Stat(backupFile); err == nil {
				fmt.Printf("💾 Taille du fichier : %.2f KB\n", float64(info.Size())/1024)
			}
			
			return
		}
	}

	fmt.Println("❌ Cette sauvegarde ne contient pas de métadonnées.")
	fmt.Println("💡 Il s'agit probablement d'une ancienne sauvegarde.")
}

func restoreBackup(backupFile, newDbName string) error {
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("le fichier de sauvegarde '%s' n'existe pas", backupFile)
	}

	newDbPath := filepath.Join("../../databases", newDbName)
	
	if _, err := os.Stat(newDbPath); err == nil {
		return fmt.Errorf("une base de données avec le nom '%s' existe déjà", newDbName)
	}

	err := os.MkdirAll(newDbPath, os.ModePerm)
	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(backupFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		path := filepath.Join(newDbPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
} 