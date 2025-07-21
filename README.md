# TVTrackr API

A modern, hexagonal-architecture backend in Go to manage users and track TV series using TheTVDB and PostgreSQL.

---

## Features
- 🔒 JWT authentication
- 📺 Search and follow TV series (TheTVDB integration)
- 👤 User registration & login
- 🗄️ PostgreSQL persistence
- 🧩 Hexagonal (ports & adapters) architecture
- 📝 OpenAPI-first design
- ♻️ Live reload for development (Air)
- 📦 Environment-based configuration
- 🚀 Database migrations with golang-migrate

---

## Quickstart

### 1. Clone & Setup
```bash
git clone https://github.com/youruser/tvtrackr-go.git
cd tvtrackr-go
cp .env.example .env
# Edit .env with your TheTVDB API key and Postgres URL
```

### 2. Run Database Migrations
See [DB_MIGRATIONS.md](./DB_MIGRATIONS.md) for details.
```bash
migrate -path ./migrations -database "$POSTGRES_URL" up
```

### 3. Start the API
```bash
go run ./cmd/main.go
```
Or with live reload (if you have [Air](https://github.com/cosmtrek/air)):
```bash
air
```

---

## API Reference
- OpenAPI spec: [`api/openapi.yaml`](./api/openapi.yaml)
- Main endpoints:
  - `POST /register` — Register user
  - `POST /login` — Login and get JWT
  - `GET /series/search?query=...` — Search series
  - `POST /series/follow` — Follow a series
  - `GET /series/followed` — List followed series

---

## Project Structure
```
├── api/                # OpenAPI spec
├── cmd/                # Entrypoint
├── internal/
│   ├── adapters/       # HTTP, TheTVDB, Postgres adapters
│   ├── core/           # Domain models & services
│   └── config/         # Centralized config
├── migrations/         # SQL migrations
├── .env.example        # Example env vars
├── DB_MIGRATIONS.md    # Migration instructions
└── README.md           # This file
```

---

## Tech Stack
- Go 1.21+
- [chi](https://github.com/go-chi/chi) (HTTP router)
- [pgx](https://github.com/jackc/pgx) (Postgres driver)
- [golang-migrate](https://github.com/golang-migrate/migrate) (DB migrations)
- [Air](https://github.com/cosmtrek/air) (live reload)
- [godotenv](https://github.com/joho/godotenv) (env vars)

---

## Contributing
Pull requests welcome! For major changes, open an issue first to discuss what you would like to change.

---

## License
MIT
# tvtrackr-go