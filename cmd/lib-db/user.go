package main

import (
	"fmt"	
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleUser(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : user <add|remove|update|grant|revoke>")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Println("Usage : user add <username> <password>")
			return
		}
		err := database.AddUser(args[1], args[2])
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}
		fmt.Println("Utilisateur ajouté.")

	case "remove":
		if len(args) < 2 {
			fmt.Println("Usage : user remove <username>")
			return
		}
		err := database.RemoveUser(args[1])
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}

	case "update":
		if len(args) < 3 {
			fmt.Println("Usage : user update <username> <new_password>")
			return
		}
		err := database.UpdateUser(args[1], args[2])
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}
		fmt.Println("Utilisateur mis à jour.")

	case "grant":
		if len(args) < 3 {
			fmt.Println("Usage : user grant <username> <dbname>")
			return
		}
		err := database.GrantDatabaseAccess(args[1], args[2])
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}
		fmt.Println("Accès accordé.")

	case "revoke":
		if len(args) < 3 {
			fmt.Println("Usage : user revoke <username> <dbname>")
			return
		}
		err := database.RevokeDatabaseAccess(args[1], args[2])
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}
		fmt.Println("Accès retiré.")

	case "reload":
		err := database.ReloadUsers()
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}

	default:
		fmt.Println("Commande inconnue :", args[0])
	}
}
