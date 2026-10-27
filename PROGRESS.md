# PROGRESS

## Current status
Phase: 7 — Documentation & diagrams   |   Last gate: PASS (2026-10-27)
Next action: Implement Phase 8 - Hardening & final audit (run final audit checklist)

## Phases
- [x] Phase 0 — Bootstrap
- [x] Phase 1 — Core infrastructure
- [x] Phase 2 — Authentication
- [x] Phase 3 — Inventory
- [x] Phase 4 — Customers
- [x] Phase 5 — Seeds
- [x] Phase 6 — Tests
- [x] Phase 7 — Documentation & diagrams
  - [x] ER_DIAGRAM.md, API.md, openapi.yaml, TESTING.md, README.md
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
