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

func LinkTables(database, table1, table2 string) error {
	reader := bufio.NewReader(os.Stdin)
	if database == "" {
		fmt.Print("Nom de la base de données : ")
		database, _ = reader.ReadString('\n')
		database = strings.TrimSpace(database)
	}
	if table1 == "" {
		fmt.Print("Nom de la table 1 : ")
		table1, _ = reader.ReadString('\n')
		table1 = strings.TrimSpace(table1)
	}
	if table2 == "" {
		fmt.Print("Nom de la table 2 : ")
		table2, _ = reader.ReadString('\n')
		table2 = strings.TrimSpace(table2)
	}

	path := fs.GetSchemaFilePath(database)
	_, err := tableExists(path, table1)
	if err != nil {
		fmt.Printf("erreur lors de la vérification de l'existence de la table %s", table1)
		return err
	}
	_, err = tableExists(path, table2)
	if err != nil {
		fmt.Printf("erreur lors de la vérification de l'existence de la table %s", table2)
		return err
	}

	fmt.Print("Type de relation ? (1:N / N:N) : ")
	relType, _ := reader.ReadString('\n')
	relType = strings.ToLower(strings.TrimSpace(relType))

	switch relType {
	case "1:n":
		fmt.Printf("Dans quelle table ajouter la clé étrangère ? (%s ou %s) : ", table1, table2)
		childTable, _ := reader.ReadString('\n')
		childTable = strings.TrimSpace(childTable)

		var parentTable string
		if childTable == table1 {
			parentTable = table2
		} else if childTable == table2 {
			parentTable = table1
		} else {
			return fmt.Errorf("table invalide")
		}

		fieldName := fmt.Sprintf("%s_id", parentTable)
		AddField(database, childTable, fieldName, "int", false)
		fmt.Printf("Relation 1:N ajoutée : %s.%s → %s.id\n", childTable, fieldName, parentTable)

	case "n:n":
		joinTable := fmt.Sprintf("%s_%s", table1, table2)
		err := AddTable(database, joinTable)
		if err != nil {
			return fmt.Errorf("échec création table de jointure : %v", err)
		}

		AddField(database, joinTable, fmt.Sprintf("%s_id", table1), "int", false)
		AddField(database, joinTable, fmt.Sprintf("%s_id", table2), "int", false)

		fmt.Printf("Relation N:N ajoutée avec la table de jointure \"%s\"\n", joinTable)

	default:
		return fmt.Errorf("relation inconnue : %s", relType)
	}

	fmt.Println("Relation ajoutée avec succès")
	return nil
}

func UnlinkTables(database, table1, table2 string) error {
	reader := bufio.NewReader(os.Stdin)
	if database == "" {
		fmt.Print("Nom de la base de données : ")
		database, _ = reader.ReadString('\n')
		database = strings.TrimSpace(database)
	}
	if table1 == "" {
		fmt.Print("Nom de la table 1 : ")
		table1, _ = reader.ReadString('\n')
		table1 = strings.TrimSpace(table1)
	}
	if table2 == "" {
		fmt.Print("Nom de la table 2 : ")
		table2, _ = reader.ReadString('\n')
		table2 = strings.TrimSpace(table2)
	}
	path := fs.GetSchemaFilePath(database)

	t1Exists, err := tableExists(path, table1)
	if err != nil {
		return err
	}
	t2Exists, err := tableExists(path, table2)
	if err != nil {
		return err
	}
	if !t1Exists || !t2Exists {
		return fmt.Errorf("une ou les deux tables n'existent pas (%s, %s)", table1, table2)
	}

	joinTable1 := fmt.Sprintf("%s_%s", table1, table2)
	joinTable2 := fmt.Sprintf("%s_%s", table2, table1)

	if exists, _ := tableExists(path, joinTable1); exists {
		return RemoveTable(database, joinTable1)
	}
	if exists, _ := tableExists(path, joinTable2); exists {
		return RemoveTable(database, joinTable2)
	}

	field1 := fmt.Sprintf("%s_id", table1)
	field2 := fmt.Sprintf("%s_id", table2)

	err1 := RemoveField(database, table1, field2)
	if err1 == nil {
		fmt.Printf("Relation supprimée : champ %s supprimé de %s\n", field2, table1)
		return nil
	}

	err2 := RemoveField(database, table2, field1)
	if err2 == nil {
		fmt.Printf("Relation supprimée : champ %s supprimé de %s\n", field1, table2)
		return nil
	}

	return fmt.Errorf("aucune relation trouvée entre %s et %s", table1, table2)
}