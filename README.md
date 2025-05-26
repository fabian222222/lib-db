# lib-db

Une librairie Go simple pour gérer des bases de données relationnelles en mémoire et sur disque.

## Fonctionnalités

- Création et gestion de bases de données
- Définition de schémas de tables avec contraintes
- Opérations CRUD (Create, Read, Update, Delete)
- Système de cache pour améliorer les performances
- Persistance des données sur disque
- Support des types de données basiques (string, int, float, bool)
- Gestion des contraintes d'unicité et de clés primaires
- Interface en ligne de commande (CLI)

## Installation

```bash
go get github.com/fabian222222/lib-db
```

## Utilisation

### Interface en ligne de commande (CLI)

Pour utiliser l'interface en ligne de commande, exécutez :

```bash
go run cmd/cli/main.go
```

Commandes disponibles :

- `help` - Affiche l'aide
- `create <db_name>` - Crée une nouvelle base de données
- `use <db_name>` - Sélectionne une base de données
- `create-table <name> <definitions>` - Crée une nouvelle table
- `insert <table> <values>` - Insère des données
- `select <table> [where conditions]` - Récupère des données
- `update <table> set values where conditions` - Met à jour des données
- `delete from <table> where conditions` - Supprime des données
- `exit` - Quitte le programme

Exemples d'utilisation de la CLI :

```bash
# Créer une nouvelle base de données
create mydb

# Sélectionner la base de données
use mydb

# Créer une table avec des colonnes
create-table users id:int:primary name:string email:string:unique age:int

# Insérer des données
insert users id=1 name=John email=john@example.com age=30

# Récupérer des données
select users
select users where name=John

# Mettre à jour des données
update users set age=31 where id=1

# Supprimer des données
delete from users where id=1
```

### Utilisation en tant que librairie

#### Créer une nouvelle base de données

```go
database, err := db.NewDatabase("ma_base")
if err != nil {
    log.Fatal(err)
}
```

#### Définir une table

```go
usersTable := schema.Table{
    Name: "users",
    Columns: []schema.Column{
        {
            Name:     "id",
            Type:     "int",
            Nullable: false,
            Unique:   true,
        },
        {
            Name:     "name",
            Type:     "string",
            Nullable: false,
        },
    },
    Primary: []string{"id"},
}

err := database.CreateTable(usersTable)
```

#### Insérer des données

```go
user := map[string]interface{}{
    "id":   1,
    "name": "John Doe",
}

err := database.Insert("users", user)
```

#### Récupérer des données

```go
// Récupérer tous les utilisateurs
users, err := database.Select("users", nil)

// Rechercher avec des conditions
conditions := map[string]interface{}{
    "name": "John Doe",
}
users, err := database.Select("users", conditions)
```

#### Mettre à jour des données

```go
conditions := map[string]interface{}{
    "name": "John Doe",
}
updates := map[string]interface{}{
    "name": "John Smith",
}
err := database.Update("users", conditions, updates)
```

#### Supprimer des données

```go
conditions := map[string]interface{}{
    "name": "John Smith",
}
err := database.Delete("users", conditions)
```

## Structure des fichiers

Chaque base de données est stockée dans un dossier sous `databases/` avec les fichiers suivants :

- `schema.txt` : Définition du schéma de la base de données
- `data.txt` : Données stockées
- `cache.txt` : Cache pour améliorer les performances

## Exemple complet

Voir le dossier `cmd/example` pour un exemple complet d'utilisation de la librairie.

## Limitations

- Pas de support pour les jointures
- Types de données limités aux types basiques
- Pas de support pour les transactions
- Pas de support pour les index
- Pas de support pour les vues

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou une pull request.
