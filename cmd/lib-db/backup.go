package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/fabian222222/lib-db/pkg/database"
)

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

	err = database.CreateBackup(dbName, backupFile, session.Username)
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
	
	canRestore, originalOwner, err := database.CanUserRestoreBackup(backupFile, session.Username)
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

	err = database.RestoreBackup(backupFile, newDbName)
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

	metadata, err := database.GetBackupInfo(backupFile)
	if err != nil {
		fmt.Println("Erreur lors de la lecture des métadonnées :", err)
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
} 