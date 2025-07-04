package main

import (
	"fmt"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleTable(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : table <add|delete|update|list>")
		return
	}

	switch args[0] {
	case "add":
		var dbName, tableName string

		if len(args) > 1 {
			dbName = args[1]
		}
		if len(args) > 2 {
			tableName = args[2]
		}
		database.AddTable(dbName, tableName)
	case "delete":
		var dbName, tableName string

		if len(args) > 1 {
			dbName = args[1]
		}
		if len(args) > 2 {
			tableName = args[2]
		}
		database.RemoveTable(dbName, tableName)
	case "update":
		var dbName, oldName, newName string

		if len(args) > 1 {
			dbName = args[1]
		}
		if len(args) > 2 {
			oldName = args[2]
		}
		if len(args) > 3 {
			newName = args[3]
		}
		database.UpdateTableName(dbName, oldName, newName)
	default:
		fmt.Printf("Commande inconnue : %s\n", args[0])
	}
}