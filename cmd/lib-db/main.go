package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Commande requise : login, logout, whoami, user, db, table, field")
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
	case "db":
		handleDb(os.Args[2:])
	case "table":
		handleTable(os.Args[2:])
	case "field":
		handleField(os.Args[2:])
	case "data":
		handleData(os.Args[2:])
	default:
		fmt.Printf("Commande inconnue : %s\n", os.Args[1])
	}
}