package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
)

type SelectQuery struct {
	DBName string            `json:"dbName"`
	Table  string            `json:"table"`
	Where  map[string]string `json:"where"`
}

type CachedSelect struct {
	Query  SelectQuery        `json:"query"`
	Result []map[string]string   `json:"result"` 
}

func SaveSelectCache(query SelectQuery, result []map[string]string) error {
	cachePath := filepath.Join("./../../databases", query.DBName, "cache.txt")

	var cache []CachedSelect

	if content, err := os.ReadFile(cachePath); err == nil && len(content) > 0 {
		json.Unmarshal(content, &cache)
	}

	for _, entry := range cache {
		if reflect.DeepEqual(entry.Query, query) {
			return nil
		}
	}

	cache = append(cache, CachedSelect{
		Query:  query,
		Result: result,
	})

	content, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, content, 0644)
}

func GetCachedSelectResult(query SelectQuery) ([]map[string]string, bool, error) {
	cachePath := filepath.Join("./../../databases", query.DBName, "cache.txt")

	content, err := os.ReadFile(cachePath)
	if err != nil || len(content) == 0 {
		return nil, false, nil
	}

	var cache []CachedSelect
	if err := json.Unmarshal(content, &cache); err != nil {
		return nil, false, err
	}

	for _, entry := range cache {
		if reflect.DeepEqual(entry.Query, query) {
			return entry.Result, true, nil
		}
	}

	return nil, false, nil
}
