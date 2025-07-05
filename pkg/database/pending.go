package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/fabian222222/lib-db/pkg/fs"
)

type CachedQuery struct {
	Action string            `json:"action"`
	DBName string            `json:"dbName"`
	Table  string            `json:"table"`
	Data   map[string]string `json:"data"`
}

func SaveQueryToCache(query CachedQuery) error {
	dirPath := filepath.Join("./../../databases", query.DBName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("database %s does not exist", query.DBName)
	}

	cachePath := filepath.Join(dirPath, "pending.txt")
	content, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, content, 0644)
}

func ClearCacheFile(dbName string) error {
	cachePath := filepath.Join("./../../databases", dbName, "pending.txt")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return fmt.Errorf("le fichier pending.txt n'existe pas pour la base de données %s", dbName)
	}

	err := os.WriteFile(cachePath, []byte{}, 0644)
	if err != nil {
		return fmt.Errorf("échec de la suppression du contenu de pending.txt : %v", err)
	}

	return nil
}

func ExecutePendingTransaction(dbName string) error {
	cachePath := filepath.Join("./../../databases", dbName, "pending.txt")
	content, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("erreur lecture pending.txt : %v", err)
	}
	if len(content) == 0 {
		return nil
	}
	
	var tx CachedQuery
	if err := json.Unmarshal(content, &tx); err != nil {
		return fmt.Errorf("format pending invalide : %v", err)
	}
	switch tx.Action {
	case "insert":
		id := tx.Data["id"]
		exists := fs.DoesDataFileExist(dbName, tx.Table, id)
		if exists {
			fmt.Println("Insertion déjà effectuée. Nettoyage du cache.")
			return ClearCacheFile(dbName)
		}
		if err := InsertData(tx.DBName, tx.Table, tx.Data); err != nil {
			return fmt.Errorf("échec insert transactionnelle : %v", err)
		}
		fmt.Println("Insertion récupérée depuis pending.txt effectuée.")
	case "update":
		id := tx.Data["id"]
		if id == "" {
			return fmt.Errorf("update invalide : id manquant dans le cache")
		}
		raw, err := os.ReadFile(fs.GetDataFile(dbName, tx.Table, id))
		if err != nil {
			return fmt.Errorf("échec de lecture avant update : %v", err)
		}
	
		var oldData map[string]string
		if err := json.Unmarshal(raw, &oldData); err != nil {
			return fmt.Errorf("données existantes mal formées : %v", err)
		}
		upToDate := true
		for k, v := range tx.Data {
			if k == "id" {
				continue
			}
			fmt.Println(oldData[k], v)
			if oldData[k] != v {
				upToDate = false
				break
			}
		}
		if upToDate {
			fmt.Println("Update déjà effectué. Nettoyage du cache.")
		}
		ClearCacheFile(dbName)
		if err := UpdateData(tx.DBName, tx.Table, id, tx.Data); err != nil {
			return fmt.Errorf("échec update transactionnelle : %v", err)
		}

		if !upToDate {
			fmt.Println("Update récupéré depuis pending.txt effectuée.")
		}
	case "delete":
		id := tx.Data["id"]
		if id == "" {
			return fmt.Errorf("delete invalide : id manquant dans le cache")
		}
		exists := fs.DoesDataFileExist(dbName, tx.Table, id)
		if !exists {
			fmt.Println("Suppression déjà effectuée. Nettoyage du cache.")
			return ClearCacheFile(dbName)
		}
		if err := DeleteData(tx.DBName, tx.Table, id); err != nil {
			return fmt.Errorf("échec delete transactionnelle : %v", err)
		}
		fmt.Println("Suppression récupérée depuis pending.txt effectuée.")
	default:
		return fmt.Errorf("action inconnue : %s", tx.Action)
	}

	return ClearCacheFile(dbName)
}