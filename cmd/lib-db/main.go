package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Commande requise : login, logout, whoami, user, reload")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "login":
		handleLogin(os.Args[2:])
	case "logout":
		handleLogout()
	case "whoami":
		handleWhoami()
	case "user":
		handleUser(os.Args[2:])
	default:
		fmt.Printf("Commande inconnue : %s\n", os.Args[1])
	}
}