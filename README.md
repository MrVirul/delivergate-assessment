# Ice Cream Shop Inventory Management System

Backend API built with Go, Fiber, GORM, PostgreSQL, JWT authentication, and RBAC.

## Features
- JWT authentication with logout (token revocation)
- Role-based access control (Owner, Manager, Cashier)
- Inventory items CRUD with search/filter/sort/pagination
- Customer records CRUD
- PostgreSQL with migrations, SQLite for tests
- Unit and integration tests with RBAC matrix

## Prerequisites
- Go 1.22+
- Docker

## Quick Start
```bash
cp .env.example .env
docker compose up -d db
make migrate
make seed
make run
```

## Testing
```bash
go test -race -count=1 ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

## Seed Accounts (dev)
| Role | Email | Password |
|------|-------|----------|
| Owner | owner@icecream.local | Owner@12345 |
| Manager | manager@icecream.local | Manager@12345 |
| Cashier | cashier@icecream.local | Cashier@12345 |

## Documentation
- [API Documentation](docs/API.md)
- [ER Diagram](docs/ER_DIAGRAM.md)
- [OpenAPI Spec](docs/openapi.yaml)
- [Testing Guide](docs/TESTING.md)
- [Decisions Log](docs/DECISIONS.md)

## Web UI
- Open http://localhost:8080/ after starting the server
- Log in with seeded accounts (see above)
- Static files in web/ served by the API
