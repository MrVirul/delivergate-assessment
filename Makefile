.PHONY: run build test test-cover lint migrate seed db-up db-down reset-db docker-build tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/seed ./cmd/seed

test:
	go test -race -count=1 ./...

test-cover:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

lint:
	gofmt -l .
	go vet ./...

migrate:
	go run ./cmd/migrate

seed:
	go run ./cmd/seed

db-up:
	docker compose up -d db

db-down:
	docker compose down

reset-db:
	docker compose down -v
	docker compose up -d db
	@sleep 3
	go run ./cmd/migrate
	go run ./cmd/seed

docker-build:
	docker build -t icecream-inventory .

tidy:
	go mod tidy
