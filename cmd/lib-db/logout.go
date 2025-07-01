package main

import (
	"fmt"
	"github.com/fabian222222/lib-db/pkg/database"
)

func handleLogout() {
	err := database.ClearSession()
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	fmt.Println("Déconnecté.")
}