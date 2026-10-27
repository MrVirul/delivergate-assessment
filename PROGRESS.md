# PROGRESS

## Current status
Phase: 0 — Bootstrap   |   Last gate: PASS (2026-10-27)
Next action: Begin Phase 1 - create core infrastructure (config, database, models, httpx)

## Phases
- [x] Phase 0 — Bootstrap
  - [x] git init, go mod
  - [x] directory skeleton
  - [x] .gitignore, .env.example
  - [x] Makefile, docker-compose.yml, Dockerfile
  - [x] PROGRESS.md, docs/DECISIONS.md
- [ ] Phase 1 — Core infrastructure
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
- go build ./... (pending)
- docker compose up -d db
