# Subscription Service (Hexagonal Architecture in Go)

A clean, modular microservice using **Hexagonal Architecture** (Ports & Adapters) with REST APIs, business logic, and a database. Designed to manage subscriptions and calculate costs over time.

---

##  Project Structure
```
subscription/
├── api
│   └── openapi.yaml # (OpenAPI spec)
├── cmd
│   └── server/
│          └── main.go # Startup, dependency wiring
├── core
│   ├── domain       # Domain entities
│   ├── ports        # Interfaces (repository, service)
│   └── usecase
├── docker-compose.prod.yml
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.dev
├── go.mod
├── go.sum
├── internal
│   └── api
│        └── genetated  # openapi generated files
│   ├── config       # Environment loading and DSN builder
│   ├── handler      # Driver adapter realization (openapi rest adapter)
│   ├── logger       # ZeroLog setup
│   └── repository   # GORM models, DB connection, repository impl
└── README.md
```

---
##  Getting Started

### Prerequisites

- Go (1.21+)
- PostgreSQL

### Technologies
- gorm
- openapi specification
- ogen generation

#### Environments naming
- Production: `.env.prod`
- Development: `.env.dev`
- Local: `.env.local`

_Example:_ [.env.example](.env.example)

Ogen command to generate OpenAPI files:
```bash
ogen --target internal/api/generated api/openapi.yaml
```

### Setup

1. **Clone the repo**  
   ```bash
   git clone https://github.com/qqax/subscriptions
   cd subscription
    ```
2. **Lifecycle**
    ```bash
    # Build development image
    docker compose --env-file .env.dev build
   
    # Start the development environment
    docker compose --env-file .env.dev up
    
    # Stop and remove containers, networks, etc.
    docker compose down
    
    # View logs
    docker compose logs -f app
    
    # Rebuild and restart containers
    docker compose up -d --build
    
    # Access the database container
    docker compose exec postgres psql -U user -d app_db
    ```


## Production
**need certificates!**
```bash
    # Build production image
    docker compose -t --env-file .env.prod build

    # Start the production environment
    docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```
