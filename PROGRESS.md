# PROGRESS

## Current status
Phase: 6 — Tests   |   Last gate: PASS (2026-10-27)
Next action: Implement Phase 7 - Documentation & diagrams (ER_DIAGRAM.md, API.md, openapi.yaml, TESTING.md, README.md)

## Phases
- [x] Phase 0 — Bootstrap
- [x] Phase 1 — Core infrastructure
- [x] Phase 2 — Authentication
- [x] Phase 3 — Inventory
- [x] Phase 4 — Customers
- [x] Phase 5 — Seeds
- [x] Phase 6 — Tests
  - [x] testutil, auth/config/httpx/service/db tests
  - [x] RBAC matrix, auth integration tests
- [ ] Phase 7 — Documentation & diagrams
- [ ] Phase 8 — Hardening & final audit

## Decisions (links to docs/DECISIONS.md)

## Blockers
None

## How to run
- go build ./... 
- go test -race -count=1 ./...
- make migrate
- make seed
- make run
- docker compose up -d db
