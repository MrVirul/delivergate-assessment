# Testing

## Running Tests
```bash
go test -race -count=1 ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

## Coverage Notes
Tests cover auth, config, httpx, database seeding, services, and integration/RBAC tests.
