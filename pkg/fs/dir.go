package fs

import (
	"os"
	"fmt"
	"path/filepath"
)

func DoesDirExist(name string) bool {
	if name == "" {
		return false
	}
	info, err := os.Stat("./../../databases/" + name)

	if err != nil {
		return false 
	}
	return info.IsDir()
}

func CreateDir(databaseName, name string) error {
	path := filepath.Join("./../../databases", databaseName, name)

	if DoesDirExist(path) {
		return fmt.Errorf("le dossier \"%s\" existe déjà", path)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du dossier \"%s\" : %w", path, err)
	}

	return nil
}