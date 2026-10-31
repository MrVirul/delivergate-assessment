# Ice Cream Shop Inventory Management System

A full-stack inventory management solution for an ice cream shop. It provides a
Go REST API for inventory items, customers, authentication, and role-based
access control, plus an optional WPF desktop client for Windows.

## Features

- JWT authentication with logout/token revocation
- Role-based access control for Owner, Manager, and Cashier
- Inventory CRUD with search, category and low-stock filters, sorting, and pagination
- Customer CRUD with search, sorting, and pagination
- Idempotent database migrations and seed data
- Optional WPF desktop client for login and inventory/customer management
- OpenAPI specification, ER diagram, and automated tests

## Stack

| Layer    | Technology                    |
| -------- | ----------------------------- |
| Backend  | Go, Fiber, GORM               |
| Database | PostgreSQL (local or Neon)    |
| Auth     | JWT (HS256), bcrypt passwords |
| Desktop  | WPF on .NET 8 (Windows only)  |

## Project Structure

| Path                           | Purpose                                                           |
| ------------------------------ | ----------------------------------------------------------------- |
| [`backend/`](backend/)         | Go Fiber REST API, migrations, and seed commands                  |
| [`desktop-app/`](desktop-app/) | Optional WPF (.NET 8) desktop client                              |
| [`docs/`](docs/)               | API reference, OpenAPI spec, and ER diagram                       |

## Prerequisites

- Go 1.27+ (see `backend/go.mod` for the exact version)
- PostgreSQL 14+ or a [Neon](https://neon.tech/) PostgreSQL database
- `make` (optional; the equivalent `go run` commands are listed below)
- For the desktop client: Windows 10/11 and the [.NET 8 SDK](https://dotnet.microsoft.com/download/dotnet/8.0)

## Configuration

All backend configuration comes from environment variables. A committed
template is available at [`backend/.env.example`](backend/.env.example); copy it
to `backend/.env` and fill in the required values:

```bash
cd backend
cp .env.example .env
```

| Variable                   | Required | Description                                                     |
| -------------------------- | -------- | --------------------------------------------------------------- |
| `DATABASE_URL`             | Yes      | PostgreSQL connection string                                    |
| `JWT_SECRET`               | Yes      | Random signing secret, at least 32 characters                   |
| `APP_ENV`                  | No       | `development` (default), `test`, or `production`                |
| `APP_PORT`                 | No       | API port (default `8080`)                                       |
| `JWT_TTL_MINUTES`          | No       | Access-token lifetime in minutes (default `60`)                 |
| `AUTO_MIGRATE`             | No       | Apply migrations automatically on server start (default `true`) |
| `BCRYPT_COST`              | No       | Password hashing cost (default `12`)                            |
| `CORS_ALLOW_ORIGINS`       | No       | Comma-separated allowed origins                                 |
| `LOGIN_RATE_LIMIT_PER_MIN` | No       | Login attempts per minute per client                            |
| `SEED_*`                   | No       | Override the default seed user emails/passwords                 |

`backend/.env` is gitignored and must never be committed. The desktop client has
its own configuration template at
[`desktop-app/appsettings.Local.example.json`](desktop-app/appsettings.Local.example.json);
see the [desktop README](desktop-app/README.md#configuration).

## Quick Start

### 1. Run the backend

From the repository root:

```bash
cd backend
cp .env.example .env
# Edit .env: set DATABASE_URL and JWT_SECRET (at least 32 characters)
make migrate
make seed
make run
```

The API starts at `http://localhost:8080` and exposes a health check at
`GET http://localhost:8080/healthz`.

If `make` is unavailable, run the same steps directly:

```bash
go run ./cmd/migrate
go run ./cmd/seed
go run ./cmd/api
```

`make run` also applies migrations automatically when `AUTO_MIGRATE=true`.

### 2. Run the optional WPF desktop client

The desktop app requires Windows and the .NET 8 SDK. Start the backend first,
then from the repository root:

```bash
dotnet run --project desktop-app
```

Or from inside the desktop project:

```bash
cd desktop-app
dotnet restore
dotnet run
```

The client connects to `http://localhost:8080/api/` by default. To point it at a
different backend, copy `desktop-app/appsettings.Local.example.json` to
`desktop-app/appsettings.Local.json` and edit `Api:BaseUrl`, or set the
`DESKTOPAPP_Api__BaseUrl` environment variable.

See [`desktop-app/README.md`](desktop-app/README.md) for desktop configuration
and permission details.

## Seed Data

`make seed` is idempotent and safe to run more than once. It creates the
following development users:

| Role    | Email              | Password      |
| ------- | ------------------ | ------------- |
| Owner   | `owner@mail.com`   | `Owner@123`   |
| Manager | `manager@mail.com` | `Manager@123` |
| Cashier | `cashier@mail.com` | `Cashier@123` |

These are the defaults from [`backend/.env.example`](backend/.env.example). If
your `backend/.env` overrides any `SEED_*` values, use those credentials instead.

It also creates sample records when `APP_ENV` is not `production`:

- 10 Sri Lankan-inspired inventory items across `flavor`, `topping`, `cone`,
  `cup`, and `packaging` categories, including a low-stock item (`Mango`)
- 5 customers with Sri Lankan names, emails, phone numbers, and addresses
  (Nimal Perera, Sanduni Fernando, Kasun Silva, Dilini Jayawardena, Chamara
  Bandara)

To use different initial users, set `SEED_OWNER_EMAIL`, `SEED_OWNER_PASSWORD`,
and the equivalent Manager/Cashier variables in `backend/.env` before the first
seed run. Seeding does not overwrite existing users, so use a fresh database (or
remove the existing user rows) to change credentials after the first run. The
three role users are created in every environment; sample items and customers
are only created outside production.

## Testing

```bash
cd backend
make test
make test-cover
```

Or directly:

```bash
go test -race -count=1 ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

The test suite covers authentication, configuration, HTTP helpers, seeding,
services, and integration/RBAC behavior.

## Documentation

- [API Documentation](docs/API.md)
- [OpenAPI Specification](docs/openapi.yaml)
- [ER Diagram](docs/ER_DIAGRAM.md)
- [Backend README](backend/README.md)
- [Desktop App README](desktop-app/README.md)

## Design Decisions

- **Layered backend**: handlers call services, which call repository interfaces,
  keeping HTTP and database concerns separate and services unit-testable.
- **PostgreSQL via GORM**: portable models/migrations, with SQLite used only for
  the in-memory test suite.
- **JWT + revocation table**: stateless access tokens while still supporting
  server-side logout.
- **RBAC enforced server-side**: the API permissions matrix is authoritative;
  the desktop client mirrors it for UI affordances.
- **Consistent API envelope**: success responses use `{"data": ...}` with list
  metadata, and errors use a stable `{"error": {...}}` shape.
- **Money stored as integer cents**: avoids floating-point rounding for prices.
- **Idempotent seeding**: repeatable setup for development and test environments.
