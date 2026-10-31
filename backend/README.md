# Ice Cream Shop Backend (Go Fiber)

RESTful API backend for the Ice Cream Shop Inventory Management System. Built
with Go, Fiber, GORM, PostgreSQL, JWT authentication, and role-based access
control.

## Prerequisites

- Go 1.27+ (see `go.mod` for the exact version)
- PostgreSQL 14+ or a Neon PostgreSQL database
- `make` (optional; direct `go run` commands are listed below)

## Configuration

Configuration is read from environment variables, with an optional `backend/.env`
file loaded automatically. Start from the committed template:

```bash
cp .env.example .env
```

Set at least:

- `DATABASE_URL` &mdash; PostgreSQL connection string
- `JWT_SECRET` &mdash; random string of at least 32 characters

Other variables (defaults shown in [`.env.example`](.env.example)) control the
port, token lifetime, automatic migrations, bcrypt cost, CORS origins, login
rate limiting, and seed credentials. `.env` is gitignored; never commit it.

## Running the Backend

```bash
cp .env.example .env
# Edit .env and set DATABASE_URL and JWT_SECRET
make migrate
make seed
make run
```

The API starts at `http://localhost:8080`, and `GET /healthz` reports health.
`make run` applies migrations automatically when `AUTO_MIGRATE=true`.

Without `make`:

```bash
go run ./cmd/migrate   # create/update tables
go run ./cmd/seed      # create users and development sample data
go run ./cmd/api       # start the HTTP server
```

### Make targets

| Target | Description |
| --- | --- |
| `make run` | Start the API server |
| `make build` | Build API, migrate, and seed binaries into `bin/` |
| `make migrate` | Run database migrations |
| `make seed` | Seed users and development data |
| `make test` | Run all tests with the race detector |
| `make test-cover` | Run tests and print coverage |
| `make lint` | Run `gofmt` and `go vet` |
| `make tidy` | Run `go mod tidy` |

## Seed Data

`make seed` (or `go run ./cmd/seed`) is idempotent and can be run repeatedly.
It creates these development accounts:

| Role | Email | Password |
| --- | --- | --- |
| Owner | `owner@mail.com` | `Owner@123` |
| Manager | `manager@mail.com` | `Manager@123` |
| Cashier | `cashier@mail.com` | `Cashier@123` |

When `APP_ENV` is not `production`, it also creates:

- 10 Sri Lankan-inspired inventory items across `flavor`, `topping`, `cone`,
  `cup`, and `packaging` categories, including a low-stock `Mango` item
- 5 sample customers with Sri Lankan names, emails, phone numbers, and addresses

User credentials can be overridden with the `SEED_OWNER_*`, `SEED_MANAGER_*`,
and `SEED_CASHIER_*` environment variables. Existing users are left unchanged,
so use a fresh database (or remove the existing user rows) if you need to change
credentials after the first seed. The three role users are seeded in every
environment; sample items and customers are skipped in production.

## Testing

```bash
make test
make test-cover
```

Equivalent commands:

```bash
go test -race -count=1 ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

Tests cover authentication, configuration, HTTP helpers, seeding, services, and
integration/RBAC behavior.

## Project Layout

```text
backend/
  cmd/api/        API entry point
  cmd/migrate/    migration entry point
  cmd/seed/       seed entry point
  internal/
    auth/         JWT, bcrypt, permission matrix
    config/       environment configuration
    database/     connection, migrations, seeders
    dto/          request/response payloads
    handler/      Fiber handlers
    httpx/        envelopes, errors, pagination, validation
    middleware/   auth, RBAC, error handling
    models/       GORM models
    repository/   data-access interfaces
    router/       route registration
    service/      business logic
  tests/          integration and RBAC tests
```

## Documentation

See the [`docs/`](../docs/) directory in the repository root:

- [API Documentation](../docs/API.md)
- [ER Diagram](../docs/ER_DIAGRAM.md)
- [OpenAPI Spec](../docs/openapi.yaml)
