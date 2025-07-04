package main

import (
	"fmt"
	"github.com/fabian222222/lib-db/pkg/database"
	"strings"
)

func handleField(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : field <add|delete|update|list>")
		return
	}

	switch args[0] {
	case "add":
		var dbName, tableName, fieldName, fieldType string
		var fieldOptionsArray []string

		if len(args) > 1 {
			dbName = args[1]
		}
		if len(args) > 2 {
			tableName = args[2]
		}
		if len(args) > 3 {
			fieldName = args[3]
		}
		if len(args) > 4 {
			fieldType = args[4]
		}
		if len(args) > 5 {
			fieldOptionsArray = strings.Split(args[5], ",")
		}

		database.AddField(dbName, tableName, fieldName, fieldType, true, fieldOptionsArray...)
	case "delete":
		var dbName, tableName, fieldName string
		if len(args) > 1 {
			dbName = args[1]
		}
		if len(args) > 2 {
			tableName = args[2]
		}
		if len(args) > 3 {
			fieldName = args[3]
		}
		database.RemoveField(dbName, tableName, fieldName)
	case "update":
		var dbName, tableName, fieldName, fieldType string
		var fieldOptionsArray []string

		if len(args) > 1 {
			dbName = args[1]
		}
		if len(args) > 2 {
			tableName = args[2]
		}
		if len(args) > 3 {
			fieldName = args[3]
		}
		if len(args) > 4 {
			fieldType = args[4]
		}
		if len(args) > 5 {
			fieldOptionsArray = strings.Split(args[5], ",")
		}
		database.UpdateField(dbName, tableName, fieldName, fieldType, fieldOptionsArray...)
	case "list":
		var dbName string
		if len(args) > 1 {
			dbName = args[1]
		}
		database.GetSchema(dbName)
	default:
		fmt.Printf("Commande inconnue : %s\n", args[0])
	}
}