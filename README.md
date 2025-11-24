# URL Shortener - Service de Raccourcissement d'URLs

Service web performant de raccourcissement et de gestion d'URLs en Go. L'application transforme une URL longue en une URL courte unique, gère la redirection instantanée et enregistre les clics de manière asynchrone. Un moniteur vérifie périodiquement la disponibilité des URLs.

## ✨ Fonctionnalités Implémentées

### 🔗 Gestion des Liens

- **Génération de codes courts uniques** : Codes de 6 caractères alphanumériques avec gestion automatique des collisions
- **Validation des URLs** : Vérification de format et détection des doublons
- **Redirection instantanée** : Redirection HTTP 302 sans latence
- **Statistiques** : Comptage des clics par lien

### 📊 Analytics Asynchrone

- **Enregistrement non-bloquant** : Workers utilisant des goroutines et channels bufferisés
- **Pool de workers configurable** : Nombre de workers ajustable via configuration
- **Traçabilité** : Capture de l'adresse IP, User-Agent et timestamp

### 🔍 Monitoring d'URLs

- **Vérification périodique** : Contrôle automatique de la disponibilité des URLs longues
- **Détection de changements d'état** : Notifications dans les logs lors de changements
- **Intervalle configurable** : Fréquence des vérifications paramétrable

### 🛠️ Architecture Technique

- **API REST** avec Gin
- **CLI complète** avec Cobra
- **Base de données SQLite** via GORM
- **Configuration** avec Viper
- **Gestion d'erreurs personnalisée** (système custom)
- **Patterns** : Repository et Service pour une architecture propre

## 🚀 Démarrage Rapide

### Prérequis

- Go 1.24+ installé
- Git

### Installation

```powershell
# Cloner le dépôt
git clone https://github.com/nghSia/UrlShortener.git
cd UrlShortener

# Télécharger les dépendances
go mod download
go mod tidy

# Compiler le binaire
go build -o url-shortener.exe
```

### Configuration

Le fichier `configs/config.yaml` permet de configurer :

- **Port du serveur** : `8080` par défaut
- **Base de données** : `url_shortener.db`
- **Analytics** : Taille du buffer (1000) et nombre de workers (5)
- **Monitor** : Intervalle de vérification (5 minutes)

## 📖 Utilisation

### 1. Initialiser la Base de Données

```powershell
.\url-shortener.exe migrate
```

### 2. Lancer le Serveur

```powershell
.\url-shortener.exe run-server
```

Le serveur démarre sur `http://localhost:8080` avec :

- API REST accessible
- Workers de clics en arrière-plan (5 workers)
- Moniteur d'URLs actif (vérification toutes les 5 min)

### 3. Commandes CLI

#### Créer un lien court

```powershell
.\url-shortener.exe create --url="https://www.google.com"
```

**Retour :**

```json
{
  "short_code": "aB3Xy9",
  "long_url": "https://www.google.com",
  "full_short_url": "http://localhost:8080/aB3Xy9"
}
```

#### Consulter les statistiques

```powershell
.\url-shortener.exe stats --code="aB3Xy9"
```

**Retour :**

```
Statistiques pour le lien 'aB3Xy9':
  URL Longue: https://www.google.com
  Clics totaux: 42
```

## 🌐 API REST

### Health Check

```powershell
curl http://localhost:8080/health
```

**Réponse :** `{"status": "ok"}`

### Créer un Lien

```powershell
curl -X POST http://localhost:8080/api/v1/links `
  -H "Content-Type: application/json" `
  -d '{"long_url": "https://example.com"}'
```

### Obtenir les Infos d'un Lien

```powershell
curl http://localhost:8080/api/v1/links/aB3Xy9
```

### Obtenir les Statistiques

```powershell
curl http://localhost:8080/api/v1/links/aB3Xy9/stats
```

### Redirection (dans le navigateur)

```
http://localhost:8080/aB3Xy9
```

→ Redirection instantanée + enregistrement asynchrone du clic

## 🧪 Tests

### Tester la création et redirection

```powershell
# 1. Créer un lien
.\url-shortener.exe create --url="https://github.com"

# 2. Tester la redirection dans le navigateur
# Ouvrir: http://localhost:8080/aB3Xy9

# 3. Vérifier les stats
.\url-shortener.exe stats --code="aB3Xy9"
```

### Tester le monitoring

```powershell
# 1. Créer plusieurs liens
.\url-shortener.exe create --url="https://google.com"
.\url-shortener.exe create --url="https://github.com"

# 2. Observer les logs du serveur
# Le moniteur vérifie automatiquement l'état toutes les 5 minutes
# Logs affichés: "[MONITOR]" et "[NOTIFICATION]" si changement d'état
```

### Tester la concurrence

```powershell
# Créer plusieurs liens simultanément pour tester les workers
for ($i=1; $i -le 10; $i++) {
  .\url-shortener.exe create --url="https://example.com/page$i"
}
```

## 🛠️ Technologies & Outils

### Frameworks & Bibliothèques

- **[Gin](https://gin-gonic.com/)** : Framework web rapide pour l'API REST
- **[Cobra](https://cobra.dev/)** : Construction de l'interface CLI
- **[Viper](https://github.com/spf13/viper)** : Gestion de configuration YAML
- **[GORM](https://gorm.io/)** : ORM pour SQLite

### Fonctionnalités Go

- **Goroutines & Channels** : Traitement asynchrone des clics
- **crypto/rand** : Génération sécurisée de codes courts
- **net/http** : Vérification de disponibilité des URLs
- **sync.Mutex** : Synchronisation de l'état du moniteur

### Patterns Architecturaux

- **Repository Pattern** : Abstraction de la couche de données
- **Service Pattern** : Logique métier centralisée
- **Worker Pool** : Traitement concurrent des événements

## 📁 Architecture du Projet

```
url-shortener/
├── cmd/                        # Points d'entrée CLI
│   ├── root.go                 # Commande racine Cobra + configuration globale
│   ├── server/
│   │   └── server.go           # Lance serveur API + workers + moniteur
│   └── cli/
│       ├── create.go           # Crée un lien court via CLI
│       ├── stats.go            # Affiche statistiques d'un lien
│       └── migrate.go          # Exécute migrations GORM
│
├── internal/                   # Code métier privé
│   ├── api/
│   │   └── handlers.go         # Handlers HTTP (routes Gin)
│   ├── models/
│   │   ├── link.go             # Modèle GORM Link
│   │   └── click.go            # Modèle GORM Click + ClickEvent
│   ├── services/
│   │   ├── link_service.go     # Génération codes + validation
│   │   └── click_service.go    # Statistiques de clics
│   ├── repository/
│   │   ├── link_repository.go  # CRUD liens (interface + GORM)
│   │   └── click_repository.go # CRUD clics (interface + GORM)
│   ├── workers/
│   │   └── click_workers.go    # Pool goroutines pour analytics async
│   ├── monitor/
│   │   └── url_monitor.go      # Surveillance périodique URLs
│   └── config/
│       └── config.go           # Structure configuration + Viper
│
├── configs/
│   └── config.yaml             # Configuration (port, DB, workers, etc.)
├── go.mod                      # Dépendances Go
└── url_shortener.db            # Base SQLite (générée automatiquement)
```

### Flux de Données

1. **Création de lien** : CLI/API → LinkService → LinkRepository → SQLite
2. **Redirection** : API Handler → LinkService → Channel → Workers → ClickRepository
3. **Monitoring** : Ticker → UrlMonitor → HTTP HEAD → Logs

## 🔧 Fonctionnalités Techniques Détaillées

### Génération de Codes Courts

- **Algorithme** : `crypto/rand` pour génération cryptographiquement sécurisée
- **Format** : 6 caractères alphanumériques (a-z, A-Z, 0-9)
- **Anti-collision** : Système de retry (max 5 tentatives)
- **Unicité** : Vérification en base avant insertion

### Analytics Asynchrone

```
Requête → Handler → Channel (buffer 1000) → Workers (pool de 5) → BDD
                 ↓
             Redirection 302 (instantanée, non bloquée)
```

**Avantages** :

- ✅ Redirection sans latence
- ✅ Traitement concurrent des clics
- ✅ Résistance aux pics de charge

### Monitoring d'URLs

- **Méthode** : Requêtes HTTP HEAD (légères)
- **Critère** : Status 2xx/3xx = accessible
- **Timeout** : 5 secondes par URL
- **État** : Map thread-safe (`sync.Mutex`)
- **Notifications** : Logs sur changement d'état

### Gestion d'Erreurs Personnalisée

Système d'erreurs custom intégré (non encore committé) :

- Erreurs typées par domaine métier
- Messages contextualisés
- Codes HTTP appropriés

## 🎯 Cas d'Usage

### 1. Service de Liens Marketing

```powershell
# Créer des liens courts pour campagnes
.\url-shortener.exe create --url="https://promo.site.com/black-friday-2025"
# → http://localhost:8080/aB3Xy9

# Tracker les clics en temps réel
.\url-shortener.exe stats --code="aB3Xy9"
```

### 2. Surveillance de Services

```powershell
# Ajouter URLs de services critiques
.\url-shortener.exe create --url="https://api.production.com/health"

# Le moniteur vérifie automatiquement et alerte sur changement
```

## 📝 Notes de Développement


### Améliorations Futures

- [ ] URLs personnalisées (custom aliases)
- [ ] Expiration automatique des liens
- [ ] Rate limiting par IP
- [ ] Dashboard web pour analytics
- [ ] Export des statistiques (CSV/JSON)

## 📄 Licence

Projet académique - TP Go Final

---

**Auteur** : Huu-Nghia TRAN, Jordy PEREIRA-ELENGA MAKOUALA, Nino FAZER, Romain MONMARCHE  
**Repository** : [github.com/nghSia/UrlShortener](https://github.com/nghSia/UrlShortener)
