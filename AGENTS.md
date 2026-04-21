# AGENTS.md

## Commands

```bash
# Run the server (requires Postgres - see docker-compose below)
go run cmd/server/main.go

# Run tests (no DB required - uses fake repositories)
go test ./...

# Build binary
go build -o server ./cmd/server
```

## Setup

1. Copy `.env.example` to `.env` and configure (or use `DATABASE_URL`)
2. Start Postgres: `docker-compose up -d`
3. Run `go run cmd/server/main.go` - server auto-migrates tables on startup

## Architecture

- **Entry point**: `cmd/server/main.go`
- **DB**: GORM with PostgreSQL, auto-migrates on startup
- **Routes**: `/health`, `/products`, `/categories`
- **Modules**: `internal/modules/{products,categories}/` - each has repository, service, handler
- **Tests**: Use `Fake*Repository` structs in same package, no external dependencies

## Notes

- Database uses `DATABASE_URL` env var if set, otherwise individual `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `SSL_MODE`, `TZ` vars
- Server auto-finds available port (3000, 3001, ..., 8080, 8081) if preferred port is in use
- No linter or formatter configured - run `go vet ./...` for basic checks
