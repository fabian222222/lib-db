package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/fabian222222/lib-db/pkg/database"
)

type DatabaseStats struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size_bytes"`
	TableCount   int       `json:"table_count"`
	LastModified time.Time `json:"last_modified"`
}

type PerformanceStats struct {
	TotalDatabases    int             `json:"total_databases"`
	TotalSizeBytes    int64           `json:"total_size_bytes"`
	DatabaseStats     []DatabaseStats `json:"database_stats"`
	GeneratedAt       time.Time       `json:"generated_at"`
	ActiveConnections int             `json:"active_connections"`
}

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
		stats, err := generatePerformanceStats()
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
		userStats, err := generateUserPerformanceStats(username)
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

	dbStats, err := generateDatabaseStats(dbName)
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
	var stats *PerformanceStats
	var err error
	
	if username == "admin" {
		stats, err = generatePerformanceStats()
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}
	} else {
		stats, err = generateUserPerformanceStats(username)
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}
	}

	statsDir := "../../stats"
	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		err = os.MkdirAll(statsDir, 0755)
		if err != nil {
			fmt.Println("Erreur lors de la création du dossier stats :", err)
			return
		}
	}
	
	userStatsDir := filepath.Join(statsDir, username)
	if _, err := os.Stat(userStatsDir); os.IsNotExist(err) {
		err = os.MkdirAll(userStatsDir, 0755)
		if err != nil {
			fmt.Println("Erreur lors de la création du dossier utilisateur :", err)
			return
		}
	}

	fileName := fmt.Sprintf("stats_export_%s.json", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(userStatsDir, fileName)
	
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Println("Erreur lors de la sérialisation :", err)
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture :", err)
		return
	}

	if username == "admin" {
		fmt.Printf("✅ Statistiques globales exportées dans : %s\n", filePath)
	} else {
		fmt.Printf("✅ Vos statistiques exportées dans : %s\n", filePath)
	}
}

func generatePerformanceStats() (*PerformanceStats, error) {
	dbPath := "../../databases"
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, err
	}

	stats := &PerformanceStats{
		GeneratedAt:       time.Now(),
		DatabaseStats:     []DatabaseStats{},
	}

	for _, file := range files {
		if file.Name() == ".session" {
			continue
		}

		dbStats, err := generateDatabaseStats(file.Name())
		if err != nil {
			continue
		}

		stats.DatabaseStats = append(stats.DatabaseStats, *dbStats)
		stats.TotalSizeBytes += dbStats.Size
		stats.TotalDatabases++
	}

	return stats, nil
}

func generateUserPerformanceStats(username string) (*PerformanceStats, error) {
	dbPath := "../../databases"
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, err
	}

	stats := &PerformanceStats{
		GeneratedAt:       time.Now(),
		DatabaseStats:     []DatabaseStats{},
	}

	for _, file := range files {
		if file.Name() == ".session" {
			continue
		}

		if !database.UserHasAccess(username, file.Name()) {
			continue
		}

		dbStats, err := generateDatabaseStats(file.Name())
		if err != nil {
			continue
		}

		stats.DatabaseStats = append(stats.DatabaseStats, *dbStats)
		stats.TotalSizeBytes += dbStats.Size
		stats.TotalDatabases++
	}

	return stats, nil
}

func generateDatabaseStats(dbName string) (*DatabaseStats, error) {
	dbPath := filepath.Join("../../databases", dbName)
	
	var totalSize int64
	var lastModified time.Time
	tableCount := 0

	err := filepath.Walk(dbPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			totalSize += info.Size()
			if info.ModTime().After(lastModified) {
				lastModified = info.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	dataPath := filepath.Join(dbPath, "data")
	if dataFiles, err := os.ReadDir(dataPath); err == nil {
		tableCount = len(dataFiles)
	}

	return &DatabaseStats{
		Name:         dbName,
		Size:         totalSize,
		TableCount:   tableCount,
		LastModified: lastModified,
	}, nil
}

func exportDatabaseStats(dbName, username string) {
	if username != "admin" && !database.UserHasAccess(username, dbName) {
		fmt.Printf("❌ Vous n'avez pas accès à la base de données '%s'\n", dbName)
		fmt.Println("💡 Seuls les propriétaires peuvent exporter les statistiques détaillées.")
		return
	}

	dbPath := filepath.Join("../../databases", dbName)
	
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("❌ La base de données '%s' n'existe pas\n", dbName)
		return
	}

	dbStats, err := generateDatabaseStats(dbName)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	statsDir := "../../stats"
	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		err = os.MkdirAll(statsDir, 0755)
		if err != nil {
			fmt.Println("Erreur lors de la création du dossier stats :", err)
			return
		}
	}
	
	dbStatsDir := filepath.Join(statsDir, dbName)
	if _, err := os.Stat(dbStatsDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbStatsDir, 0755)
		if err != nil {
			fmt.Println("Erreur lors de la création du dossier pour la base :", err)
			return
		}
	}

	fileName := fmt.Sprintf("stats_%s_%s.json", dbName, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(dbStatsDir, fileName)
	
	data, err := json.MarshalIndent(dbStats, "", "  ")
	if err != nil {
		fmt.Println("Erreur lors de la sérialisation :", err)
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture :", err)
		return
	}

	fmt.Printf("✅ Statistiques de la base '%s' exportées dans : %s\n", dbName, filePath)
}