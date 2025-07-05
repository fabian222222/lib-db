package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"github.com/fabian222222/lib-db/pkg/fs"
	"github.com/lucsky/cuid"
	"path/filepath"
	"io/ioutil"
)

func InsertData(databaseName, tableName string, rawInputs ...map[string]string) error {
	reader := bufio.NewReader(os.Stdin)

	if databaseName == "" {
		fmt.Print("Nom de la base de données : ")
		databaseName, _ = reader.ReadString('\n')
		databaseName = strings.TrimSpace(databaseName)
	}

	if !fs.DoesDirExist(databaseName) {
		return fmt.Errorf("la base de données \"%s\" n'existe pas", databaseName)
	}

	if tableName == "" {
		fmt.Print("Nom de la table : ")
		tableName, _ = reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)
	}

	schema, err := GetSchema(databaseName)
	if err != nil {
		return err
	}

	fields, ok := schema[tableName]
	if !ok {
		return fmt.Errorf("la table \"%s\" n'existe pas", tableName)
	}

	tableDir := filepath.Join("./../../databases", databaseName, "data", tableName)
	if !fs.DoesDirExist(tableDir) {
		err := os.MkdirAll(tableDir, 0755)
		if err != nil {
			return fmt.Errorf("échec création dossier pour table : %w", err)
		}
	}

	var fieldNames []string
	for _, field := range fields {
		parts := strings.Split(field, ":")
		fieldName := strings.TrimSpace(parts[0])
		fieldNames = append(fieldNames, fieldName)
	}

	for _, input := range rawInputs {
		entry := map[string]string{}
		entry["id"] = cuid.New()
		fmt.Printf("ID généré : %s\n", entry["id"])

		for _, field := range fieldNames {
			if field == "id" {
				continue
			}
			val, ok := input[field]
			if !ok || strings.TrimSpace(val) == "" {
				fmt.Printf("Valeur pour \"%s\" : ", field)
				val, _ = reader.ReadString('\n')
				val = strings.TrimSpace(val)
			}
			if strings.HasSuffix(field, "_id") && val != "" {
				relatedTable := strings.TrimSuffix(field, "_id")
				relatedTablePath := filepath.Join(databaseName, "data", relatedTable)
				if !fs.DoesDirExist(relatedTablePath) {
					return fmt.Errorf("la table liée \"%s\" n'existe pas pour la clé étrangère \"%s\"", relatedTable, field)
				}

				if !fs.DoesFileExist(filepath.Join("./../../databases", databaseName, "data", relatedTable, val+".json")) {
					return fmt.Errorf("la valeur \"%s\" pour \"%s\" n'existe pas dans la table \"%s\"", val, field, relatedTable)
				}
			}
			entry[field] = val
		}

		data, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			return err
		}

		SaveQueryToCache(CachedQuery{
			Action: "insert",
			DBName: databaseName,
			Table:  tableName,
			Data:   entry,
		})
		filename := filepath.Join(tableDir, entry["id"]+".json")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func UpdateData(databaseName, tableName, targetID string, updates map[string]string) error {
	reader := bufio.NewReader(os.Stdin)

	if databaseName == "" {
		fmt.Print("Nom de la base de données : ")
		databaseName, _ = reader.ReadString('\n')
		databaseName = strings.TrimSpace(databaseName)
	}
	if !fs.DoesDirExist(databaseName) {
		return fmt.Errorf("La base de données \"%s\" n'existe pas", databaseName)
	}

	if tableName == "" {
		fmt.Print("Nom de la table : ")
		tableName, _ = reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)
	}

	if targetID == "" {
		fmt.Print("ID de la donnée à modifier : ")
		targetID, _ = reader.ReadString('\n')
		targetID = strings.TrimSpace(targetID)
	}

	schema, err := GetSchema(databaseName)
	if err != nil {
		return err
	}
	fields, ok := schema[tableName]
	if !ok {
		return fmt.Errorf("La table \"%s\" n'existe pas", tableName)
	}

	validFields := make(map[string]bool)
	for _, f := range fields {
		parts := strings.Split(f, ":")
		validFields[strings.TrimSpace(parts[0])] = true
	}

	dataFile := fs.GetDataFile(databaseName, tableName, targetID)
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return fmt.Errorf("L'entrée avec ID \"%s\" n'existe pas dans la table \"%s\"", targetID, tableName)
	}

	fileBytes, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	var entry map[string]string
	if err := json.Unmarshal(fileBytes, &entry); err != nil {
		return fmt.Errorf("Erreur de lecture du fichier JSON : %v", err)
	}

	for field := range validFields {
		if field == "id" {
			continue
		}

		val, ok := updates[field]
		if !ok {
			fmt.Printf("Nouvelle valeur pour \"%s\" (laissez vide pour conserver \"%s\") : ", field, entry[field])
			input, _ := reader.ReadString('\n')
			val = strings.TrimSpace(input)
		}

		if val == "" {
			continue
		}

		if strings.HasSuffix(field, "_id") {
			relatedTable := strings.TrimSuffix(field, "_id")
			relatedTablePath := filepath.Join(databaseName, "data", relatedTable)

			if !fs.DoesDirExist(relatedTablePath) {
				return fmt.Errorf("La table liée \"%s\" n'existe pas pour la clé étrangère \"%s\"", relatedTable, field)
			}

			if !fs.DoesFileExist(filepath.Join("./../../databases", databaseName, "data", relatedTable, val+".json")) {
				return fmt.Errorf("la valeur \"%s\" pour \"%s\" n'existe pas dans la table \"%s\"", val, field, relatedTable)
			}
		}
		entry[field] = val
	}

	updatedJSON, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}


	SaveQueryToCache(CachedQuery{
		Action: "update",
		DBName: databaseName,
		Table:  tableName,
		Data:   entry,
	})

	if err := os.WriteFile(dataFile, updatedJSON, 0644); err != nil {
		return fmt.Errorf("Erreur lors de l'écriture : %v", err)
	}

	fmt.Println("Entrée mise à jour avec succès.")
	return nil
}

func DeleteData(databaseName, tableName, id string) error {
	if databaseName == "" {
		return fmt.Errorf("le nom de la base de données ne peut pas être vide")
	}
	if tableName == "" {
		return fmt.Errorf("le nom de la table ne peut pas être vide")
	}
	if id == "" {
		return fmt.Errorf("l'id ne peut pas être vide")
	}

	if !fs.DoesDirExist(databaseName) {
		return fmt.Errorf("la base de données \"%s\" n'existe pas", databaseName)
	}

	filePath := fs.GetDataFile(databaseName, tableName, id)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("l'entrée avec l'id \"%s\" n'existe pas dans la table \"%s\"", id, tableName)
	}

	SaveQueryToCache(CachedQuery{
		Action: "delete",
		DBName: databaseName,
		Table:  tableName,
		Data:   map[string]string{"id": id},
	})

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de l'entrée : %v", err)
	}

	fmt.Printf("Entrée avec l'id \"%s\" supprimée avec succès.\n", id)
	return nil
}

func SelectData(databaseName, tableName string, whereClauses map[string]string) ([]map[string]string, error) {
	if databaseName == "" {
		return nil, fmt.Errorf("le nom de la base de données ne peut pas être vide")
	}
	if tableName == "" {
		return nil, fmt.Errorf("le nom de la table ne peut pas être vide")
	}

	if !fs.DoesDirExist(databaseName) {
		return nil, fmt.Errorf("la base de données \"%s\" n'existe pas", databaseName)
	}

	query := SelectQuery{
		DBName: databaseName,
		Table:  tableName,
		Where:  whereClauses,
	}

	cachedResults, found, err := GetCachedSelectResult(query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du cache : %v", err)
	}
	if found {
		fmt.Println("Résultat récupéré depuis le cache.")
		return cachedResults, nil
	}

	tablePath := fs.GetDataFilePath(databaseName, tableName)

	matchingEntries := []map[string]string{}

	files, err := ioutil.ReadDir(tablePath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le dossier table: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(tablePath, file.Name())
		dataBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("impossible de lire le fichier %s: %v", file.Name(), err)
		}

		var entry map[string]string
		err = json.Unmarshal(dataBytes, &entry)
		if err != nil {
			return nil, fmt.Errorf("erreur d'unmarshal JSON dans %s: %v", file.Name(), err)
		}

		match := true
		for k, v := range whereClauses {
			val, ok := entry[k]
			if !ok || val != v {
				match = false
				break
			}
		}

		if match {
			matchingEntries = append(matchingEntries, entry)
		}
	}

	SaveSelectCache(query, matchingEntries)

	return matchingEntries, nil
}
