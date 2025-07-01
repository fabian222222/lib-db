package main

import (
	"fmt"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleLogin(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage : login <username> <password>")
		return
	}

	username := args[0]
	password := args[1]

	ok, _, err := database.Authenticate(username, password)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	if !ok {
		fmt.Println("Identifiants incorrects")
		return
	}

	err = database.SaveSession(username)
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde de session :", err)
		return
	}

	fmt.Printf("Connecté en tant que %s\n", username)
}