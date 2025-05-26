package database

import(
	"fmt"
	"os"
	"github.com/fabian222222/lib-db/pkg/fs"
)

func CreateDatabase(name string) {
	isDirExist := fs.DoesDirExist(name)
	dbPath := "./databases/" + name
	if isDirExist {
		fmt.Println("Database", name, "already exist")
		return
	}

	err := os.MkdirAll(dbPath, os.ModePerm)

	if err != nil {
        fmt.Println("Error creating database:", err)
        return
    }

    fmt.Println("Database", name, "created at", dbPath)
	fs.CreateFile(name, "schema.txt")
	fs.CreateFile(name, "data.txt")
	fs.CreateFile(name, "cache.txt")
	return
}