package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"log" 
	"bufio"
	"strings"
	"path/filepath"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Databases []string `json:"databases"`
}

type Session struct {
	Username string `json:"username"`
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

const usersFilePath = "../../users/users.json"
const sessionFilePath = "../../databases/.session"

func SaveSession(username string) error {
	data, err := json.Marshal(Session{Username: username})
	if err != nil {
		return err
	}
	return os.WriteFile(sessionFilePath, data, 0644)
}

func LoadSession() (*Session, error) {
	data, err := os.ReadFile(sessionFilePath)
	if err != nil {
		return nil, err
	}
	var s Session
	err = json.Unmarshal(data, &s)
	return &s, err
}

func ClearSession() error {
	return os.Remove(sessionFilePath)
}

func Authenticate(username, password string) (bool, *User, error) {
	users, err := LoadUsers()
	if err != nil {
		return false, nil, err
	}

	for _, u := range users {
		if u.Username == username {
			if CheckPasswordHash(password, u.Password) {
				SaveSession(username)
				return true, &u, nil
			}
			return false, nil, nil
		}
	}

	return false, nil, nil
}

func IsAuthenticated() (bool, *Session, error) {
	session, err := LoadSession()
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, session, nil
}

func UserHasAccess(username, dbName string) bool {
	users, err := LoadUsers()
	if err != nil {
		return false
	}
	for _, u := range users {
		if u.Username == username {
			for _, db := range u.Databases {
				if db == dbName {
					return true
				}
			}
		}
	}
	return false
}

func ensureUsersFile() error {
	dir := filepath.Dir(usersFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(usersFilePath); os.IsNotExist(err) {
		f, err := os.Create(usersFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write([]byte("[]"))
		return err
	}
	return nil
}

func LoadUsers() ([]User, error) {
	if err := ensureUsersFile(); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(usersFilePath)
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(data, &users)
	return users, err
}

func SaveUsers(users []User) error {
	ok, _, err := IsAuthenticated()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Vous devez être connecté pour faire cette action")
	}
	if err := ensureUsersFile(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(usersFilePath, data, 0644)
}


func AddUser(username, password string) error {
	ok, _, err := IsAuthenticated()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Vous devez être connecté pour faire cette action")
	}
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Username == username {
			return errors.New("username already exists")
		}
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := User{
		Username:  username,
		Password:  string(hashedPwd),
		Databases: []string{},
	}

	users = append(users, newUser)
	return SaveUsers(users)
}

func UpdateUser(username, newPassword string) error {
	ok, _, err := IsAuthenticated()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Vous devez être connecté pour faire cette action")
	}
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Username == username {
			if newPassword != "" {
				hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
				if err != nil {
					return err
				}
				users[i].Password = string(hashedPwd)
			}
			return SaveUsers(users)
		}
	}

	return errors.New("user not found")
}

func RemoveUser(username string) error {
	ok, _, err := IsAuthenticated()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Vous devez être connecté pour faire cette action")
	}

	users, err := LoadUsers()
	if err != nil {
		return err
	}

	found := false
	for _, u := range users {
		if u.Username == username {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Utilisateur \"%s\" non trouvé.\n", username)
		return nil
	}

	if len(users) == 1 && users[0].Username == username {
		fmt.Printf("Attention, vous êtes sur le point de supprimer le dernier utilisateur \"%s\".\n", username)
		fmt.Print("Confirmez-vous cette suppression ? (oui/non) : ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "oui" && response != "o" {
			fmt.Println("Suppression annulée.")
			return nil
		}
	}

	filtered := make([]User, 0)
	for _, u := range users {
		if u.Username != username {
			filtered = append(filtered, u)
		}
	}

	if err := SaveUsers(filtered); err != nil {
		return err
	}

	session, err := LoadSession()
	if err != nil {
		return err
	}

	if session.Username == username {
		if err := ClearSession(); err != nil {
			return fmt.Errorf("utilisateur supprimé mais erreur lors de la suppression de la session : %v", err)
		}
		fmt.Println("Votre session a été fermée car vous avez supprimé votre propre compte.")
	}

	fmt.Printf("Utilisateur \"%s\" supprimé.\n", username)
	return nil
}

func GrantDatabaseAccess(username, dbName string) error {
	ok, _, err := IsAuthenticated()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Vous devez être connecté pour faire cette action")
	}
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Username == username {
			for _, db := range u.Databases {
				if db == dbName {
					return nil
				}
			}
			users[i].Databases = append(users[i].Databases, dbName)
			return SaveUsers(users)
		}
	}

	return errors.New("user not found")
}

func RevokeDatabaseAccess(username, dbName string) error {
	ok, _, err := IsAuthenticated()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Vous devez être connecté pour faire cette action")
	}
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Username == username {
			filtered := make([]string, 0)
			for _, db := range u.Databases {
				if db != dbName {
					filtered = append(filtered, db)
				}
			}
			users[i].Databases = filtered
			return SaveUsers(users)
		}
	}

	return errors.New("user not found")
}

func ReloadUsers() error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	defaultUser := User{
		Username:  "admin",
		Password:  string(hashedPwd),
		Databases: []string{},
	}

	users := []User{defaultUser}

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println("Utilisateurs réinitialisés. Utilisateur par défaut créé : admin / admin")
	return os.WriteFile(usersFilePath, data, 0644)
}