# Ice Cream Shop Backend (Go Fiber)

RESTful API backend for the Ice Cream Shop Inventory Management System built with Go, Fiber, GORM, PostgreSQL, JWT authentication, and RBAC.

## Prerequisites
- Go 1.22+
- PostgreSQL database (e.g., Neon cloud Postgres)

## Quick Start
```bash
cp .env.example .env
# Set DATABASE_URL and other variables in .env
make migrate
make seed
make run
```

The API will start at `http://localhost:8080`.

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
See the [`docs/`](../docs/) directory in the repository root:
- [API Documentation](../docs/API.md)
- [ER Diagram](../docs/ER_DIAGRAM.md)
- [OpenAPI Spec](../docs/openapi.yaml)
- [Testing Guide](../docs/TESTING.md)
- [Decisions Log](../docs/DECISIONS.md)
