package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

func GeneratePerformanceStats() (*PerformanceStats, error) {
	dbPath := "../../databases"
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, err
	}

	stats := &PerformanceStats{
		TotalDatabases:    0,
		TotalSizeBytes:    0,
		DatabaseStats:     []DatabaseStats{},
		GeneratedAt:       time.Now(),
		ActiveConnections: 1,
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != ".session" {
			dbStats, err := GenerateDatabaseStats(file.Name())
			if err != nil {
				continue
			}
			stats.DatabaseStats = append(stats.DatabaseStats, *dbStats)
			stats.TotalDatabases++
			stats.TotalSizeBytes += dbStats.Size
		}
	}

	return stats, nil
}

func GenerateUserPerformanceStats(username string) (*PerformanceStats, error) {
	dbPath := "../../databases"
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, err
	}

	stats := &PerformanceStats{
		TotalDatabases:    0,
		TotalSizeBytes:    0,
		DatabaseStats:     []DatabaseStats{},
		GeneratedAt:       time.Now(),
		ActiveConnections: 1,
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != ".session" {
			if UserHasAccess(username, file.Name()) {
				dbStats, err := GenerateDatabaseStats(file.Name())
				if err != nil {
					continue
				}
				stats.DatabaseStats = append(stats.DatabaseStats, *dbStats)
				stats.TotalDatabases++
				stats.TotalSizeBytes += dbStats.Size
			}
		}
	}

	return stats, nil
}

func GenerateDatabaseStats(dbName string) (*DatabaseStats, error) {
	dbPath := filepath.Join("../../databases", dbName)
	
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("la base de données '%s' n'existe pas", dbName)
	}

	var totalSize int64
	var lastModified time.Time
	var tableCount int

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

func ExportStats(username string, statsDir string) (string, error) {
	var stats *PerformanceStats
	var err error
	
	if username == "admin" {
		stats, err = GeneratePerformanceStats()
		if err != nil {
			return "", err
		}
	} else {
		stats, err = GenerateUserPerformanceStats(username)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		err = os.MkdirAll(statsDir, 0755)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la création du dossier stats : %v", err)
		}
	}
	
	userStatsDir := filepath.Join(statsDir, username)
	if _, err := os.Stat(userStatsDir); os.IsNotExist(err) {
		err = os.MkdirAll(userStatsDir, 0755)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la création du dossier utilisateur : %v", err)
		}
	}

	fileName := fmt.Sprintf("stats_export_%s.json", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(userStatsDir, fileName)
	
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erreur lors de la sérialisation : %v", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'écriture : %v", err)
	}

	return filePath, nil
}

func ExportDatabaseStats(dbName, username string, statsDir string) (string, error) {
	if username != "admin" && !UserHasAccess(username, dbName) {
		return "", fmt.Errorf("vous n'avez pas accès à la base de données '%s'", dbName)
	}

	dbPath := filepath.Join("../../databases", dbName)
	
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("la base de données '%s' n'existe pas", dbName)
	}

	dbStats, err := GenerateDatabaseStats(dbName)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		err = os.MkdirAll(statsDir, 0755)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la création du dossier stats : %v", err)
		}
	}
	
	dbStatsDir := filepath.Join(statsDir, dbName)
	if _, err := os.Stat(dbStatsDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbStatsDir, 0755)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la création du dossier pour la base : %v", err)
		}
	}

	fileName := fmt.Sprintf("stats_%s_%s.json", dbName, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(dbStatsDir, fileName)
	
	data, err := json.MarshalIndent(dbStats, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erreur lors de la sérialisation : %v", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'écriture : %v", err)
	}

	return filePath, nil
} 