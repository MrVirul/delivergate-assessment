FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/seed ./cmd/seed

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /bin/api /bin/migrate /bin/seed /bin/
COPY .env.example /app/.env.example
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
