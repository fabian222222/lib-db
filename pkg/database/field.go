package database

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/fabian222222/lib-db/pkg/fs"
)

var allowedTypes = map[string]bool{
	"int":      true,
	"string":   true,
	"float":    true,
	"bool":     true,
	"datetime": true,
}

var allowedOptions = map[string]bool{
	"pk":     true,
	"unique": true,
	"fk":     true,
}

func AddField(databaseName, tableName, fieldName, fieldType string, showLogs bool, options ...string) {
	reader := bufio.NewReader(os.Stdin)

	if databaseName == "" {
		fmt.Print("Nom de la base de données : ")
		databaseName, _ = reader.ReadString('\n')
		databaseName = strings.TrimSpace(databaseName)
	}
	if tableName == "" {
		fmt.Print("Nom de la table : ")
		tableName, _ = reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)
	}
	if fieldName == "" {
		fmt.Print("Nom du champ : ")
		fieldName, _ = reader.ReadString('\n')
		fieldName = strings.TrimSpace(fieldName)
	}
	if fieldType == "" {
		fmt.Print("Type du champ (int, string, float, bool, datetime) : ")
		fieldType, _ = reader.ReadString('\n')
		fieldType = strings.TrimSpace(fieldType)
	}
	if len(options) == 0 && showLogs {
		fmt.Print("Options du champ (séparées par des virgules, ex: pk,unique ou fk=table.id) : ")
		optLine, _ := reader.ReadString('\n')
		optLine = strings.TrimSpace(optLine)
		if optLine != "" {
			options = strings.Split(optLine, ",")
		}
	}
	
	fieldDefinition := fieldName + ":" + fieldType
	if len(options) > 0 {
		fieldDefinition += ":" + strings.Join(options, ",")
	}
	
	if err := ValidateFieldDefinition(fieldDefinition); err != nil {
		fmt.Println("erreur lors de la validation du champ", err)
		return
	}
	
	path := fs.GetSchemaFilePath(databaseName)
	lines, err := fs.ReadLines(path)
	if err != nil {
		fmt.Println("la base de données n'existe pas")
		return
	}
	found := false
	newLines := []string{}
	inTable := false
	fieldExists := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)

		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			inTable = trim == "["+tableName+"]"
			if inTable {
				found = true
			}
			newLines = append(newLines, line)
			continue
		}

		if inTable {
			if strings.HasPrefix(trim, fieldName+":") {
				fieldExists = true
			}
		}

		newLines = append(newLines, line)
	}

	if !found {
		fmt.Printf("table \"%s\" introuvable dans la base \"%s\"\n", tableName, databaseName)
		return
	}
	if fieldExists {
		fmt.Printf("le champ \"%s\" existe déjà dans la table \"%s\"\n", fieldName, tableName)
		return
	}

	finalLines := []string{}
	inTable = false
	for _, line := range newLines {
		finalLines = append(finalLines, line)
		if strings.TrimSpace(line) == "["+tableName+"]" {
			inTable = true
			continue
		}
		if inTable && (line == "" || strings.HasPrefix(strings.TrimSpace(line), "[")) {
			finalLines = append(finalLines[:len(finalLines)-1], fieldDefinition, line) 
			inTable = false
		}
	}

	if inTable {
		finalLines = append(finalLines, fieldDefinition)
	}

	err = fs.WriteLines(path, finalLines)
	if err != nil {
		fmt.Println("erreur lors de l'écriture du fichier", err)
		return
	}
	if showLogs {
		fmt.Printf("le champ \"%s\" a été ajouté à la table \"%s\"\n", fieldName, tableName)
	}
	return
}

func RemoveField(database string, tableName string, fieldName string, showLogs ...bool) error {
	log := true
	if len(showLogs) > 0 {
		log = showLogs[0]
	}
	reader := bufio.NewReader(os.Stdin)

	if database == "" {
		fmt.Print("Nom de la base de données : ")
		database, _ = reader.ReadString('\n')
		database = strings.TrimSpace(database)
	}
	if tableName == "" {
		fmt.Print("Nom de la table : ")
		tableName, _ = reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)
	}
	if fieldName == "" {
		fmt.Print("Nom du champ : ")
		fieldName, _ = reader.ReadString('\n')
		fieldName = strings.TrimSpace(fieldName)
	}

	if fieldName == "id" {
		fmt.Printf("le champ \"%s\" n'est pas supprimable", fieldName)
		return nil
	}

	path := fs.GetSchemaFilePath(database)

	lines, err := fs.ReadLines(path)
	if err != nil {
		return err
	}

	fieldExists := false
	inTable := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			inTable = trim == "["+tableName+"]"
			continue
		}
		if inTable && strings.HasPrefix(trim, fieldName+":") {
			fieldExists = true
			break
		}
	}

	if !inTable {
		fmt.Printf("la table \"%s\" n'existe pas dans la base \"%s\"", tableName, database)
		return nil
	}

	if !fieldExists {
		fmt.Printf("le champ \"%s\" n'existe pas dans la table \"%s\"", fieldName, tableName)
		return nil
	}

	if log {
		fmt.Printf("Êtes-vous sûr de vouloir supprimer le champ \"%s\" de la table \"%s\" ? (oui/non) : ", fieldName, tableName)
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(strings.ToLower(confirmation))
		if confirmation != "oui" && confirmation != "o" {
			fmt.Println("Suppression annulée.")
			return nil
		}
	}

	newLines := []string{}
	inTable = false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			inTable = trim == "["+tableName+"]"
			newLines = append(newLines, line)
			continue
		}

		if inTable && strings.HasPrefix(trim, fieldName+":") {
			continue
		}
		newLines = append(newLines, line)
	}

	return fs.WriteLines(path, newLines)
}


func UpdateField(databaseName, tableName, fieldName, newType string, newOptions ...string) error {
	reader := bufio.NewReader(os.Stdin)

	if databaseName == "" {
		fmt.Print("Nom de la base de données : ")
		databaseName, _ = reader.ReadString('\n')
		databaseName = strings.TrimSpace(databaseName)
	}
	if tableName == "" {
		fmt.Print("Nom de la table : ")
		tableName, _ = reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)
	}
	if fieldName == "" {
		fmt.Print("Nom du champ à mettre à jour : ")
		fieldName, _ = reader.ReadString('\n')
		fieldName = strings.TrimSpace(fieldName)
	}
	if newType == "" {
		fmt.Print("Nouveau type du champ (int, string, float, bool, datetime) : ")
		newType, _ = reader.ReadString('\n')
		newType = strings.TrimSpace(newType)
	}
	if len(newOptions) == 0 {
		fmt.Print("Nouvelles options (séparées par des virgules, ex: pk,unique ou fk=table.id) : ")
		optLine, _ := reader.ReadString('\n')
		optLine = strings.TrimSpace(optLine)
		if optLine != "" {
			newOptions = strings.Split(optLine, ",")
		}
	}

	newDefinition := fieldName + ":" + newType
	if len(newOptions) > 0 {
		newDefinition += ":" + strings.Join(newOptions, ",")
	}

	path := fs.GetSchemaFilePath(databaseName)

	lines, err := fs.ReadLines(path)
	if err != nil {
		fmt.Println("La base de données n'existe pas")
		return nil
	}

	schema := make(map[string][]string)
	var currentTable string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			currentTable = strings.Trim(trim, "[]")
			schema[currentTable] = []string{}
		} else if currentTable != "" {
			schema[currentTable] = append(schema[currentTable], trim)
		}
	}

	fields, ok := schema[tableName]
	if !ok {
		fmt.Printf("table \"%s\" introuvable\n", tableName)
		return nil
	}

	found := false
	for _, f := range fields {
		if strings.HasPrefix(f, fieldName+":") {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("le champ \"%s\" n'existe pas dans la table \"%s\"\n", fieldName, tableName)
		return nil
	}

	if err := ValidateFieldDefinition(newDefinition); err != nil {
		fmt.Println("erreur lors de la validation du champ", err)
		return nil
	}

	if err := RemoveField(databaseName, tableName, fieldName, false); err != nil {
		fmt.Println("erreur lors de la suppression du champ", err)
		return nil
	}
	AddField(databaseName, tableName, fieldName, newType, false, newOptions...)
	fmt.Printf("le champ \"%s\" a été mis à jour dans la table \"%s\"\n", fieldName, tableName)
	return nil
}

func GetSchema(database string) (map[string][]string, error) {
	reader := bufio.NewReader(os.Stdin)

	if database == "" {
		fmt.Print("Nom de la base de données : ")
		database, _ = reader.ReadString('\n')
		database = strings.TrimSpace(database)
	}
	path := fs.GetSchemaFilePath(database)

	lines, err := fs.ReadLines(path)
	if err != nil {
		fmt.Println("La base de données n'existe pas")
		return nil, err
	}

	schema := make(map[string][]string)
	var currentTable string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			currentTable = strings.Trim(trim, "[]")
			schema[currentTable] = []string{}
		} else if currentTable != "" {
			schema[currentTable] = append(schema[currentTable], trim)
		}
	}

	fmt.Println("📘 Schéma de la base de données :", database)
	fmt.Println("──────────────────────────────────────")
	for table, fields := range schema {
		fmt.Printf("📂 Table: %s\n", table)
		for _, field := range fields {
			fmt.Printf("   └─ %s\n", field)
		}
		fmt.Println()
	}

	return schema, nil
}

func ValidateFieldDefinition(def string) error {
	parts := strings.Split(def, ":")
	if len(parts) < 2 {
		return fmt.Errorf("le champ '%s' est invalide (format attendu: nom:type:options...)", def)
	}

	fieldType := parts[1]
	if !allowedTypes[fieldType] {
		return fmt.Errorf("type non autorisé: '%s'", fieldType)
	}

	if len(parts) > 2 {
		options := strings.Split(parts[2], ",")
		for _, opt := range options {
			if strings.HasPrefix(opt, "fk=") {
				continue
			}
			if !allowedOptions[opt] {
				return fmt.Errorf("option non autorisée: '%s'", opt)
			}
		}
	}
	return nil
}