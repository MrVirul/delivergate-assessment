# Ice Cream Shop Inventory Management System

Full-stack inventory management solution comprising a Go Fiber REST API and a .NET 8 WPF desktop application.

## Project Structure
- [`backend/`](backend/) &mdash; Go Fiber REST API with PostgreSQL, JWT authentication, and RBAC.
- [`desktop-app/`](desktop-app/) &mdash; WPF (.NET 8) desktop client.
- [`docs/`](docs/) &mdash; API specification, ER diagrams, testing guide, and decisions log.

## Quick Start

### 1. Backend (Go Fiber)
```bash
cd backend
cp .env.example .env
# Set DATABASE_URL in .env to your PostgreSQL/Neon connection string
make migrate
make seed
make run
```
The API server starts on `http://localhost:8080`.

### 2. Desktop Application (WPF .NET 8)
```bash
cd desktop-app
dotnet restore
dotnet run
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
