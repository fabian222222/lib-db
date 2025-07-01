package main

import (
	"fmt"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleWhoami() {
	ok, session, err := database.IsAuthenticated()
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	if !ok {
		fmt.Println("Aucun utilisateur connecté.")
		return
	}
	fmt.Printf("Connecté en tant que %s\n", session.Username)
}