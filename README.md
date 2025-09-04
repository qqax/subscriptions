# Subscription Service (Hexagonal Architecture in Go)

A clean, modular microservice using **Hexagonal Architecture** (Ports & Adapters) with REST APIs, business logic, and a database. Designed to manage subscriptions and calculate costs over time.

---

##  Features

- CRUD operations on subscriptions including service name, monthly price, user ID (UUID), start/end dates.
- Calculate total subscription cost for a period, with optional filters by user or service.
- Structured as hexagonal architecture:
    - `domain`: pure business entities
    - `ports`: interfaces
    - `app`: use-case logic
    - `adapters`: HTTP handlers, DB connections (GORM), logger, OpenAPI spec, generated models

---

##  Project Structure
```
subscription/
├── cmd/
│ └── server/
│ └── main.go # Startup, dependency wiring
├── internal/
│ ├── domain/ # Domain entities
│ ├── ports/ # Interfaces (repository, service)
│ ├── app/ # Business logic / use-cases
│ └── adapters/
│ ├── http/ # REST handlers, OpenAPI glue
│ └── db/ # GORM models, DB connection, repository impl
├── pkg/
│ └── logger/ # ZeroLog setup
├── config/ # Environment loading and DSN builder
├── openapi/ # subscriptions.yaml (OpenAPI spec)
├── generated/ # GORM Gen generated tables
├── migrations/ # SQL migration files
├── go.mod # Module definitions
└── README.md # This file
```

---
##  Getting Started

### Prerequisites

- Go (1.21+)
- PostgreSQL

### Environments naming
- Production: `.env.prod`
- Development: `.env.dev`
- Local: `.env.local`

Example: [.env.example](.env.example)

### Setup

1. **Clone the repo**  
   ```bash
   git clone <repo-url>
   cd subscription
    ```
   
gorm
opanapi
ogen ogen --target internal/api/generated api/openapi.yaml
betteralign -apply ./...


# Простота использования через Makefile
make dev        # Запуск для разработки
make test       # Запуск тестов
make deploy     # Деплой в прод

make build      # Собрать образ
make test       # Запустить тесты
make up         # Запустить композ
make deploy     # Деплой


```bash
# Запуск development окружения
docker-compose up -d

# Остановка
docker-compose down

# Просмотр логов
docker-compose logs -f app

# Пересборка
docker-compose up -d --build

# Зайти в контейнер с БД
docker-compose exec postgres psql -U user -d app_db
```

## Production
```bash
# Сборка production образа
docker build -t my-app:prod .

# Запуск production контейнера
docker run -d \
-p 8080:8080 \
--name my-app \
--env-file .env.local \
my-app:prod
```


