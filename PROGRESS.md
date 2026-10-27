# PROGRESS

## Current status
Phase: 1 — Core infrastructure   |   Last gate: NOT RUN
Next action: Add auth package (permissions, jwt, password), models/migrations done partially, continue to phase 2 or complete phase 1 pieces

## Phases
- [x] Phase 0 — Bootstrap
  - [x] git init, go mod
  - [x] directory skeleton
  - [x] .gitignore, .env.example
  - [x] Makefile, docker-compose.yml, Dockerfile
  - [x] PROGRESS.md, docs/DECISIONS.md
- [x] Phase 1 — Core infrastructure (partial)
  - [x] config loader
  - [x] database (Postgres + SQLite openers)
  - [x] models (user, item, customer, revoked_token)
  - [x] migrate
  - [x] httpx (envelope, errors, pagination, validator)
  - [x] central error handler middleware
  - [x] cmd/migrate
- [ ] Phase 2 — Authentication
- [ ] Phase 3 — Inventory
- [ ] Phase 4 — Customers
- [ ] Phase 5 — Seeds
- [ ] Phase 6 — Tests
- [ ] Phase 7 — Documentation & diagrams
- [ ] Phase 8 — Hardening & final audit

## Decisions (links to docs/DECISIONS.md)

## Blockers
None

## How to run
- go build ./... 
- make migrate (needs DB)
- docker compose up -d db
