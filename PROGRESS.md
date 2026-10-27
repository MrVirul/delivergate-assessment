# PROGRESS

## Current status
Phase: 5 — Seeds   |   Last gate: PASS (2026-10-27)
Next action: Implement Phase 6 - Tests (testutil, auth/config/httpx tests, service tests, RBAC matrix, API tests)

## Phases
- [x] Phase 0 — Bootstrap
- [x] Phase 1 — Core infrastructure
- [x] Phase 2 — Authentication
- [x] Phase 3 — Inventory
- [x] Phase 4 — Customers
- [x] Phase 5 — Seeds
  - [x] cmd/seed, internal/database/seed.go (idempotent)
- [ ] Phase 6 — Tests
- [ ] Phase 7 — Documentation & diagrams
- [ ] Phase 8 — Hardening & final audit

## Decisions (links to docs/DECISIONS.md)

## Blockers
None

## How to run
- go build ./... 
- make migrate
- make seed
- make run
- docker compose up -d db
