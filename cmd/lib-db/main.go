package main

import (
	"fmt"
	"os"
	"github.com/fabian222222/lib-db/pkg/database"
)

func main() {
	if len(os.Args) < 2 {
        fmt.Println("Usage: lib-db <dbname>")
        return
    }
    dbName := os.Args[1]
    database.CreateDatabase(dbName)
}