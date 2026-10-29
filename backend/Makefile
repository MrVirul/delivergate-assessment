.PHONY: run build test test-cover lint migrate seed tidy

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

tidy:
	go mod tidy
