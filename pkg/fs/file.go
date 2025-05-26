package fs
import (
	"fmt"
	"os"
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
	path := "./databases/" + databaseName + "/" + fileName
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