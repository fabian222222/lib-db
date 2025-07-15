package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleStats(args []string) {
	ok, session, err := database.IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour voir les statistiques.")
		return
	}

	if len(args) == 0 {
		showGeneralStats(session.Username)
		return
	}

	switch args[0] {
	case "db":
		if len(args) < 2 {
			fmt.Println("Usage : stats db <database_name> [export]")
			return
		}
		if len(args) >= 3 && args[2] == "export" {
			exportDatabaseStats(args[1], session.Username)
		} else {
			showDatabaseStats(args[1], session.Username)
		}
	case "export":
		exportStats(session.Username)
	case "performance":
		showPerformanceReport(session.Username)
	default:
		fmt.Println("Options disponibles : db <name>, export, performance")
	}
}

func showGeneralStats(username string) {
	isAdmin := (username == "admin")
	
	if isAdmin {
		stats, err := database.GeneratePerformanceStats()
		if err != nil {
			fmt.Println("Erreur lors de la génération des statistiques :", err)
			return
		}

		fmt.Println("\n📊 STATISTIQUES GÉNÉRALES LIB-DB (MODE ADMIN)")
		fmt.Println("═══════════════════════════════════════════════")
		fmt.Printf("🗄️  Nombre total de bases de données : %d\n", stats.TotalDatabases)
		fmt.Printf("💾 Taille totale des données : %.2f MB\n", float64(stats.TotalSizeBytes)/(1024*1024))
		fmt.Printf("🔗 Connexions actives : %d\n", stats.ActiveConnections)
		fmt.Printf("⏰ Généré le : %s\n\n", stats.GeneratedAt.Format("02/01/2006 15:04:05"))

		fmt.Println("📋 DÉTAIL PAR BASE DE DONNÉES :")
		fmt.Println("─────────────────────────────────")
		for _, dbStat := range stats.DatabaseStats {
			fmt.Printf("• %s : %.2f KB (%d tables) - Modifié : %s\n", 
				dbStat.Name, 
				float64(dbStat.Size)/1024,
				dbStat.TableCount,
				dbStat.LastModified.Format("02/01 15:04"))
		}
	} else {
		userStats, err := database.GenerateUserPerformanceStats(username)
		if err != nil {
			fmt.Println("Erreur lors de la génération des statistiques :", err)
			return
		}

		fmt.Printf("\n📊 VOS STATISTIQUES LIB-DB (%s)\n", username)
		fmt.Println("═════════════════════════════════════════")
		fmt.Printf("🗄️  Vos bases de données : %d\n", len(userStats.DatabaseStats))
		fmt.Printf("💾 Taille totale de vos données : %.2f MB\n", float64(userStats.TotalSizeBytes)/(1024*1024))
		fmt.Printf("⏰ Généré le : %s\n\n", userStats.GeneratedAt.Format("02/01/2006 15:04:05"))

		if len(userStats.DatabaseStats) == 0 {
			fmt.Println("❌ Vous n'avez accès à aucune base de données.")
			fmt.Println("💡 Demandez à l'administrateur de vous accorder des permissions.")
			return
		}

		fmt.Println("📋 VOS BASES DE DONNÉES :")
		fmt.Println("────────────────────────")
		for _, dbStat := range userStats.DatabaseStats {
			fmt.Printf("• %s : %.2f KB (%d tables) - Modifié : %s\n", 
				dbStat.Name, 
				float64(dbStat.Size)/1024,
				dbStat.TableCount,
				dbStat.LastModified.Format("02/01 15:04"))
		}
	}
}

func showDatabaseStats(dbName, username string) {
	if username != "admin" && !database.UserHasAccess(username, dbName) {
		fmt.Printf("❌ Vous n'avez pas accès à la base de données '%s'\n", dbName)
		fmt.Println("💡 Seuls les propriétaires peuvent voir les statistiques détaillées.")
		return
	}

	dbPath := filepath.Join("../../databases", dbName)
	
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("❌ La base de données '%s' n'existe pas\n", dbName)
		return
	}

	dbStats, err := database.GenerateDatabaseStats(dbName)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Printf("\n📊 STATISTIQUES DÉTAILLÉES - %s\n", dbName)
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("💾 Taille totale : %.2f KB\n", float64(dbStats.Size)/1024)
	fmt.Printf("📋 Nombre de tables : %d\n", dbStats.TableCount)
	fmt.Printf("⏰ Dernière modification : %s\n", dbStats.LastModified.Format("02/01/2006 15:04:05"))

	files := []string{"schema.txt", "cache.txt", "pending.txt"}
	fmt.Println("\n📁 ANALYSE DES FICHIERS :")
	fmt.Println("─────────────────────────")
	
	for _, fileName := range files {
		filePath := filepath.Join(dbPath, fileName)
		if info, err := os.Stat(filePath); err == nil {
			fmt.Printf("• %s : %.2f KB\n", fileName, float64(info.Size())/1024)
		}
	}

	dataPath := filepath.Join(dbPath, "data")
	if info, err := os.Stat(dataPath); err == nil && info.IsDir() {
		dataSize := int64(0)
		filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				dataSize += info.Size()
			}
			return nil
		})
		fmt.Printf("• data/ : %.2f KB\n", float64(dataSize)/1024)
	}
}

func showPerformanceReport(username string) {
	isAdmin := (username == "admin")
	
	if isAdmin {
		fmt.Println("\n⚡ RAPPORT DE PERFORMANCE SYSTÈME (MODE ADMIN)")
		fmt.Println("═══════════════════════════════════════════════")
	} else {
		fmt.Printf("\n⚡ RAPPORT DE PERFORMANCE - %s\n", username)
		fmt.Println("═════════════════════════════════════════")
	}
	
	start := time.Now()
	dbPath := "../../databases"
	files, err := os.ReadDir(dbPath)
	elapsed := time.Since(start)
	
	if err != nil {
		fmt.Println("Erreur lors du test de performance :", err)
		return
	}

	fmt.Printf("🔍 Temps de listage des bases de données : %v\n", elapsed)

	fmt.Println("\n💡 RECOMMANDATIONS :")
	fmt.Println("─────────────────────")
	
	if isAdmin {
		if len(files) > 10 {
			fmt.Println("• Considérer l'archivage des anciennes bases de données")
		}
		fmt.Println("• Surveiller l'activité des utilisateurs")
		fmt.Println("• Gérer les permissions d'accès")
	} else {
		fmt.Println("• Effectuer des sauvegardes régulières de vos bases")
		fmt.Println("• Optimiser vos requêtes si nécessaire")
		fmt.Println("• Nettoyer les données inutiles")
	}
	
	fmt.Println("• Monitorer l'espace disque disponible")
}

func exportStats(username string) {
	filePath, err := database.ExportStats(username, "../../stats")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	if username == "admin" {
		fmt.Printf("✅ Statistiques globales exportées dans : %s\n", filePath)
	} else {
		fmt.Printf("✅ Vos statistiques exportées dans : %s\n", filePath)
	}
}



func exportDatabaseStats(dbName, username string) {
	filePath, err := database.ExportDatabaseStats(dbName, username, "../../stats")
	if err != nil {
		fmt.Printf("❌ %s\n", err.Error())
		return
	}

	fmt.Printf("✅ Statistiques de la base '%s' exportées dans : %s\n", dbName, filePath)
}