package database

import(
	"fmt"
	"os"
	"github.com/fabian222222/lib-db/pkg/fs"
	"bufio"
	"strings"
	"path/filepath"
)

func CreateDatabase(name string) {
	ok, _, err := IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour créer une base de données.")
		return
	}

	isDirExist := fs.DoesDirExist(name)
	dbPath := "./../../databases/" + name
	if isDirExist {
		fmt.Println("Database", name, "already exist")
		return
	}

	err = os.MkdirAll(dbPath, os.ModePerm)

	if err != nil {
        fmt.Println("Error creating database:", err)
        return
    }

    fmt.Println("Database", name, "created at", dbPath)
	fs.CreateFile(name, "schema.txt")
	fs.CreateFile(name, "data.txt")
	fs.CreateFile(name, "cache.txt")
	fs.CreateFile(name, "pending.txt")
	return
}

func UpdateDatabaseName(oldName, newName string) {
	ok, _, err := IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return 
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour renommer une base de données.")
		return 
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Voulez-vous vraiment renommer la base de données '%s' en '%s' ? (oui/non) : ", oldName, newName)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'entrée utilisateur :", err)
		return
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "oui" && response != "o" {
		fmt.Println("Renommage annulé.")
		return 
	}

	oldPath := filepath.Join("./../../databases", oldName)
	newPath := filepath.Join("./../../databases", newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		fmt.Println("La base de données", oldName, "n'existe pas")
		return
	}

	if _, err := os.Stat(newPath); err == nil {
		fmt.Println("Une base de données avec le nom", newName, "existe déjà")
		return
	}

	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("Erreur lors du renommage :", err)
		return
	}

	fmt.Printf("Base de données renommée de '%s' en '%s'.\n", oldName, newName)
	return 
}

func DeleteDatabase(name string) {
	ok, _, err := IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour supprimer une base de données.")
		return
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Êtes-vous sûr de vouloir supprimer définitivement la base de données '%s' ? (oui/non) : ", name)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'entrée utilisateur :", err)
		return
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "oui" && response != "o" {
		fmt.Println("Suppression annulée.")
		return
	}

	dbPath := filepath.Join("./../../databases", name)

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("La base de données", name, "n'existe pas")
		return
	}

	err = os.RemoveAll(dbPath)
	if err != nil {
		fmt.Println("Erreur lors de la suppression :", err)
		return
	}

	fmt.Printf("Base de données '%s' supprimée.\n", name)
	return
}

func ListDatabases() {
	ok, _, err := IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur lors de la vérification de la session :", err)
		return
	}
	if !ok {
		fmt.Println("Vous devez être connecté pour lister les bases de données.")
		return	
	}
	dbPath := "./../../databases"
	files, err := os.ReadDir(dbPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture des bases de données :", err)
		return
	}

	fmt.Println("Liste des bases de données :")
	for _, file := range files {
		if file.Name() != ".session" {
			fmt.Println(file.Name())
		}
	}
	return
}