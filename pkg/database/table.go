package database

import (
	"fmt"
	"os"
	"strings"
	"github.com/fabian222222/lib-db/pkg/fs"
	"bufio"
	"path/filepath"
)

func AddTable(database, tableName string) error {
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

	path := fs.GetSchemaFilePath(database)

	exists, err := tableExists(path, tableName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("la table \"%s\" existe déjà", tableName)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("\n[%s]\n", tableName))
	if err != nil {
		return err
	}

	AddField(database, tableName, "id", "int", true, "pk", "unique")
	fs.CreateDir(database, "data/" + tableName)
	fmt.Printf("la table \"%s\" a été créée", tableName)
	return nil
}

func UpdateTableName(database, oldTableName, newTableName string) error {
	reader := bufio.NewReader(os.Stdin)
	if database == "" {
		fmt.Print("Nom de la base de données : ")
		database, _ = reader.ReadString('\n')
		database = strings.TrimSpace(database)
	}
	if oldTableName == "" {
		fmt.Print("Ancien nom de la table : ")
		oldTableName, _ = reader.ReadString('\n')
		oldTableName = strings.TrimSpace(oldTableName)
	}
	if newTableName == "" {
		fmt.Print("Nouveau nom de la table : ")
		newTableName, _ = reader.ReadString('\n')
		newTableName = strings.TrimSpace(newTableName)
	}

	path := fs.GetSchemaFilePath(database)
	lines, err := fs.ReadLines(path)
	if err != nil {
		return err
	}

	oldFound := false
	newExists := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "["+oldTableName+"]" {
			oldFound = true
		}
		if strings.TrimSpace(line) == "["+newTableName+"]" {
			newExists = true
		}
	}
	if !oldFound {
		return fmt.Errorf("la table \"%s\" n'existe pas", oldTableName)
	}
	if newExists {
		return fmt.Errorf("la table \"%s\" existe déjà", newTableName)
	}

	newLines := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "["+oldTableName+"]" {
			newLines = append(newLines, "["+newTableName+"]")
		} else {
			newLines = append(newLines, line)
		}
	}

	oldPath := filepath.Join("./../../databases", database, oldTableName)
	newPath := filepath.Join("./../../databases", database, newTableName)

	os.Rename(oldPath, newPath)

	return fs.WriteLines(path, newLines)
}

func RemoveTable(database, tableName string) error {
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

	path := fs.GetSchemaFilePath(database)

	lines, err := fs.ReadLines(path)
	if err != nil {
		return err
	}

	newLines := []string{}
	inTable := false
	found := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			if trim == "["+tableName+"]" {
				inTable = true
				found = true
				continue
			} else {
				inTable = false
			}
		}
		if inTable {
			continue
		}
		newLines = append(newLines, line)
	}

	if !found {
		fmt.Printf("la table \"%s\" n'existe pas", tableName)
		return nil
	}

	err = fs.WriteLines(path, newLines)
	if err != nil {
		fmt.Printf("erreur lors de la suppression de la table \"%s\"", tableName)
		return err
	}

	if err := os.RemoveAll(fs.GetDataFilePath(database, tableName)); err != nil {
		return fmt.Errorf("échec de la suppression du dossier \"%s\": %w", path, err)
	} 
	fmt.Printf("la table \"%s\" a été supprimée", tableName)
	return nil
}

func tableExists(path, tableName string) (bool, error) {
	lines, err := fs.ReadLines(path)
	if err != nil {
		return false, err
	}
	for _, line := range lines {
		if strings.TrimSpace(line) == "["+tableName+"]" {
			return true, nil
		}
	}
	return false, nil
}