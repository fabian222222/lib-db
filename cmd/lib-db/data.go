package main

import (
	"fmt"
	"strings"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleData(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : data <insert|update|delete|select|cache> <database> <table> <field1=value1 field2=value2 ...>")
		return
	}

	action := args[0]
	switch action {
	case "insert":
		if len(args) < 3 {
			fmt.Println("Usage : data insert <database> <table> <field1=value1 field2=value2 ...>")
			return
		}
		var input map[string]string = map[string]string{}
		if len(args) > 3 {
			for i := 3; i < len(args); i++ {
				pair := args[i]
				if parts := strings.SplitN(pair, "=", 2); len(parts) == 2 {
					input[parts[0]] = parts[1]
				}
			}
		}
		err := database.InsertData(args[1], args[2], input)
		if err != nil {
			fmt.Println("Erreur :", err)
		}
	case "update":
		if len(args) < 4 {
			fmt.Println("Usage : data update <database> <table> <id> <field1=value1 field2=value2 ...>")
			return
		}
		var input map[string]string = map[string]string{}
		if len(args) > 2 {
			for i := 2; i < len(args); i++ {
				pair := args[i]
				if parts := strings.SplitN(pair, "=", 2); len(parts) == 2 {
					input[parts[0]] = parts[1]
				}
			}
		}
		err := database.UpdateData(args[1], args[2], args[3], input)
		if err != nil {
			fmt.Println("Erreur :", err)
		}
	case "delete":
		if len(args) < 4 {
			fmt.Println("Usage : data delete <database> <table> <id>")
			return
		}
		err := database.DeleteData(args[1], args[2], args[3])
		if err != nil {
			fmt.Println("Erreur :", err)
		}
	case "select":
		if len(args) < 3 {
			fmt.Println("Usage : data select <database> <table> [field=value ...]")
			return
		}
		databaseName := args[1]
		tableName := args[2]
		filters := make(map[string]string)
		for _, pair := range args[3:] {
			if parts := strings.SplitN(pair, "=", 2); len(parts) == 2 {
				filters[parts[0]] = parts[1]
			}
		}

		results, err := database.SelectData(databaseName, tableName, filters)
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("Aucune donnée trouvée.")
			return
		}

		for _, entry := range results {
			fmt.Println(entry)
		}
	case "cache":
		if len(args) < 2 {
			fmt.Println("Usage : data cache <database>")
			return
		}
		database.ExecutePendingTransaction(args[1])
	default:
		fmt.Println("Action non reconnue.")
	}
}
