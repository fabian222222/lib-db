# 🗄️ Lib-DB - Système de Gestion de Base de Données (POSTGRE SQL - REDIS - Transaction save)

**Un SGBD (Système de Gestion de Base de Données) personnalisé développé en Go**

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

### 1. **Interface en Ligne de Commande (CLI)**
```bash
# Authentification
./lib-db login <username> <password>
./lib-db logout
./lib-db whoami

# Gestion des utilisateurs
./lib-db user add <username> <password>
./lib-db user grant <username> <database>

# Gestion des bases de données
./lib-db db create <name>
./lib-db db list
./lib-db db delete <name>

# Gestion des tables
./lib-db table add <db> <table>
./lib-db table link <db> <table1> <table2>

# Manipulation des données
./lib-db data insert <db> <table> field1=value1 field2=value2
./lib-db data select <db> <table> [filters...]

# Fonctionnalités avancées
./lib-db backup <database> [file]      # Sauvegarde ZIP sécurisée (propriétaires uniquement)
./lib-db backup info <file>            # Voir les infos d'une sauvegarde
./lib-db restore <file> <new_db>       # Restauration (propriétaires uniquement)
./lib-db stats [db|export|performance] # Statistiques et monitoring
```

### 2. **Système de Sauvegarde Sécurisé** 💾
- **Compression ZIP** automatique avec métadonnées de propriétaire
- **Horodatage automatique** des sauvegardes
- **Permissions strictes** : seuls les propriétaires peuvent sauvegarder/restaurer
- **Vérification d'intégrité** et traçabilité complète
- **Info des sauvegardes** sans les ouvrir

### 3. **Statistiques de Performance Sécurisées** 📊
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

---

*Développé dans le cadre d'un projet académique pour démontrer la maîtrise des concepts de bases de données et du développement en Go.*