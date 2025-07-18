# 🗄️ Lib-DB - Système de Gestion de Base de Données (POSTGRE SQL - REDIS - Transaction save)

**Un SGBD (Système de Gestion de Base de Données) personnalisé développé en Go**

Contexte : Étant donné que le projet de base de données était trop léger, nous avons décidé de refaire un système pgsql en ajoutant une partie REDIS et une partie
sauvegarde de transactions.

## Membres

- ZUO Fabian
- ASSO'O EMANE Ulysse
- BOUAKI Arthur

## 📋 Vue d'ensemble

Lib-DB est un système de gestion de base de données léger et modulaire qui offre :

- 🔐 **Authentification sécurisée** avec hachage bcrypt
- 🗂️ **Gestion complète des bases de données** (CRUD)
- 💾 **Système de sauvegarde/restauration** automatique
- 📈 **Statistiques de performance** détaillées avec analyse
- 🔗 **Gestion des relations** entre tables
- 🚀 **Architecture modulaire** et extensible

## 🏗️ Architecture

```
lib-db/
├── cmd/lib-db/          # Interface CLI et serveur web
│   ├── main.go          # Point d'entrée
│   ├── login.go         # Authentification
│   ├── database.go      # Gestion des bases de données
│   ├── table.go         # Gestion des tables
│   ├── field.go         # Gestion des champs
│   ├── data.go          # Manipulation des données
│   ├── web.go           # Interface web
│   ├── backup.go        # Sauvegarde/Restauration
│   └── stats.go         # Statistiques de performance
├── pkg/
│   ├── database/        # Logique métier
│   └── fs/             # Gestion du système de fichiers
├── databases/          # Stockage des données
├── stats/              # Stockage des exports de statistiques
└── backups/          # Stockage des sauvegardes
```

## 🚀 Fonctionnalités

### 1. **Interface en Ligne de Commande (CLI) - Commandes Complètes**

#### **Authentification**

```bash
./lib-db login <username> <password>     # Connexion utilisateur
./lib-db logout                          # Déconnexion
./lib-db whoami                          # Voir l'utilisateur connecté
```

#### **Gestion des utilisateurs**

```bash
./lib-db user add <username> <password>      # Ajouter un utilisateur
./lib-db user remove <username>              # Supprimer un utilisateur
./lib-db user update <username> <new_password> # Mettre à jour le mot de passe
./lib-db user grant <username> <database>    # Accorder l'accès à une base
./lib-db user revoke <username> <database>   # Retirer l'accès à une base
./lib-db user reload                         # Recharger les utilisateurs
```

#### **Gestion des bases de données**

```bash
./lib-db db create <name>                    # Créer une base de données
./lib-db db delete <name>                    # Supprimer une base de données
./lib-db db update <old_name> <new_name>     # Renommer une base de données
./lib-db db list                             # Lister les bases de données
```

#### **Gestion des tables**

```bash
./lib-db table add <db> <table>              # Ajouter une table
./lib-db table delete <db> <table>           # Supprimer une table
./lib-db table update <db> <old_name> <new_name> # Renommer une table
./lib-db table link <db> <table1> <table2>   # Lier deux tables
./lib-db table unlink <db> <table1> <table2> # Délier deux tables
```

#### **Gestion des champs**

```bash
./lib-db field add <db> <table> <field> <type> [options]    # Ajouter un champ
./lib-db field delete <db> <table> <field>                  # Supprimer un champ
./lib-db field update <db> <table> <field> <type> [options] # Modifier un champ
./lib-db field list <db>                                    # Lister le schéma
```

#### **Manipulation des données**

```bash
./lib-db data insert <db> <table> field1=value1 field2=value2 # Insérer des données
./lib-db data update <db> <table> <id> field1=value1         # Mettre à jour
./lib-db data delete <db> <table> <id>                       # Supprimer
./lib-db data select <db> <table> [field=value ...]          # Sélectionner avec filtres
./lib-db data cache <db>                                     # Exécuter les transactions en attente
```

#### **Sauvegarde et restauration**

```bash
./lib-db backup <database> [file]           # Sauvegarde ZIP sécurisée
./lib-db backup info <file>                 # Informations d'une sauvegarde
./lib-db restore <file> <new_db>            # Restaurer une sauvegarde
```

#### **Statistiques et monitoring**

```bash
./lib-db stats                              # Statistiques générales
./lib-db stats db <name>                    # Statistiques d'une base
./lib-db stats db <name> export             # Exporter les stats d'une base
./lib-db stats export                       # Exporter toutes les statistiques
./lib-db stats performance                  # Rapport de performance
```

### 2. **User Flow - Parcours Utilisateur Complet**

#### **Scénario 1 : Flux utilisateur standard avec gestion d'erreurs**

```bash

cd cmd/lib-db

# 1. Rechargement des utilisateurs (initialisation)
go run *.go user reload

# 2. Déconnexion (si une session existe)
go run *.go logout

# 3. Tentative de création de base sans être connecté (va échouer)
go run *.go db create ecommerce
# ❌ Erreur : Vous devez être connecté

# 4. Connexion en tant qu'admin
go run *.go login admin admin

# 5. Création d'une base de données
go run *.go db create ecommerce

# 6. Création de deux tables
go run *.go table add ecommerce users
go run *.go table add ecommerce products

# 7. Création d'une liaison entre les tables
go run *.go table link ecommerce users products

# 8. Ajout de champs dans la table users
go run *.go field add ecommerce users name string
go run *.go field add ecommerce users email string unique

# 9. Ajout de champs dans la table products
go run *.go field add ecommerce products name string
go run *.go field add ecommerce products price int

# 10. Insertion de données dans users
go run *.go data insert ecommerce users name=Jean email=jean@example.com
go run *.go data insert ecommerce users name=Marie email=marie@example.com

# 11. Insertion de données dans products
go run *.go data insert ecommerce products name=Laptop price=999
go run *.go data insert ecommerce products name=Mouse price=29
go run *.go data insert ecommerce products name=Keyboard price=79

# 12. Récupération des données avec filtres (WHERE)
go run *.go data select ecommerce users name=Jean
go run *.go data select ecommerce products

# 13. Mise à jour des données
go run *.go data update ecommerce products <id> price=10000

# 14. Suppression de données
go run *.go data delete ecommerce products <id>
```

#### **Scénario 2 : Test du cache (performance)**

```bash
# 1. Premier appel (va créer le cache)
go run *.go data select ecommerce users

# 2. Deuxième appel identique (va utiliser le cache)
go run *.go data select ecommerce users

# 3. Vérification des performances
go run *.go stats performance

# 4. Exécution manuelle du cache (si nécessaire)
go run *.go data cache ecommerce
```

#### **Scénario 3 : Gestion avancée avec sauvegarde**

```bash
# 1. Voir les statistiques détaillées
go run *.go stats db ecommerce

# 2. Faire une sauvegarde sécurisée
go run *.go backup ecommerce

# 3. Voir les informations de la sauvegarde
go run *.go backup info

# 4. Exporter les statistiques
go run *.go stats export

# 5. Créer un nouvel utilisateur et lui donner accès
go run *.go user add developer password123
go run *.go user grant developer ecommerce

# 6. Voir qui est connecté
go run *.go whoami
```

### 3. **Système de Sauvegarde Sécurisé** 💾

- **Compression ZIP** automatique avec métadonnées de propriétaire
- **Horodatage automatique** des sauvegardes
- **Permissions strictes** : seuls les propriétaires peuvent sauvegarder/restaurer
- **Vérification d'intégrité** et traçabilité complète
- **Info des sauvegardes** sans les ouvrir

### 4. **Statistiques de Performance Sécurisées** 📊

- **Permissions granulaires** : utilisateurs voient seulement leurs bases
- **Mode admin** : accès aux statistiques globales du système
- **Monitoring en temps réel** avec recommandations personnalisées
- **Score de fragmentation** et analyse des performances
- **Export JSON** filtré selon les permissions

## 🔐 Sécurité

- **Hachage bcrypt** pour les mots de passe
- **Système de sessions** sécurisées
- **Permissions granulaires** par base de données
- **Authentification obligatoire** pour toutes les opérations
- **Sauvegardes sécurisées** : métadonnées de propriétaire intégrées
- **Contrôle d'accès strict** : seuls les propriétaires peuvent backup/restore
- **Statistiques filtérées** : chaque utilisateur voit seulement ses données
- **Mode administrateur** : accès complet aux statistiques système

## 💻 Installation et Utilisation

```bash
# Clone et compilation
git clone https://github.com/fabian222222/lib-db
cd lib-db/cmd/lib-db
go build -o lib-db *.go

# Initialiser les utilisateurs
./lib-db user reload

# Connexion
./lib-db login admin admin

# Création d'une base de données
./lib-db db create ecommerce

# Voir les statistiques de performance
./lib-db stats

# Faire une sauvegarde sécurisée
./lib-db backup ecommerce

# Voir les infos d'une sauvegarde
./lib-db backup info backup_ecommerce_20240115_143025.zip

# Statistiques (selon vos permissions)
./lib-db stats                    # Vos bases si utilisateur normal, tout si admin
./lib-db stats db ecommerce       # Seulement si vous y avez accès
```

## 📈 Cas d'Usage

1. **Prototypage rapide** d'applications
2. **Systèmes embarqués** nécessitant une base légère
3. **Apprentissage** des concepts de SGBD
4. **Tests et développement** sans infrastructure lourde
5. **Stockage local** pour applications desktop

## 🎯 Avantages Techniques

- **Performance** : Stockage fichier optimisé
- **Simplicité** : API intuitive et documentation claire
- **Extensibilité** : Architecture modulaire
- **Monitoring** : Outils d'analyse intégrés

## 🔧 Technologies Utilisées

- **Go 1.24.1** - Langage principal
- **bcrypt** - Sécurité des mots de passe
- **JSON** - Sérialisation des données
- **ZIP** - Compression des sauvegardes
