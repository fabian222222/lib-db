package fs
import (
	"fmt"
	"os"
	"bufio"
	"path/filepath"
)

func DoesFileExist(name string) bool {
	if name == "" {
		return false
	}
	_, err := os.Open(name)
	return err == nil
}

func CreateFile(databaseName string, fileName string) error {
	if DoesFileExist(fileName) {
		return fmt.Errorf("file %s already exists", fileName)
	}
	path := "./../../databases/" + databaseName + "/" + fileName
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("file arealdy exist", path)
	}
	if !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func GetSchemaFilePath(database string) string {
	return filepath.Join("./../../databases", database, "schema.txt")
}

func GetTableFilePath(database string) string {
	return filepath.Join("./../../databases", database, "tables.txt")
}

func GetDataFilePath(database string, tableName string) string {
	return filepath.Join("./../../databases", database, "data", tableName)
}

func GetDataFile(database string, tableName string, id string) string {
	return filepath.Join("./../../databases", database, "data", tableName, id + ".json")
}

func GetCacheFilePath(database string) string {
	return filepath.Join("./../../databases", database, "cache.txt")
}

func GetPendingFilePath(database string) string {
	return filepath.Join("./../../databases", database, "pending.txt")
}

func DoesDataFileExist(database string, tableName string, id string) bool {
	return DoesFileExist(GetDataFile(database, tableName, id))
}

func ReadLines(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}