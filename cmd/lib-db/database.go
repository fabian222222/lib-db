package main

import (
	"fmt"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleDb(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : db <create|delete|update|list>")
		return
	}

	switch args[0] {
	case "create":
		if len(args) < 2 {
			fmt.Println("Usage : db create <database_name>")
			return
		}
		database.CreateDatabase(args[1])
	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage : db delete <database_name>")
			return
		}
		database.DeleteDatabase(args[1])
	case "update":
		if len(args) < 3 {
			fmt.Println("Usage : db update <old_name> <new_name>")
			return
		}
		database.UpdateDatabaseName(args[1], args[2])
	case "list":
		database.ListDatabases()
	}
}