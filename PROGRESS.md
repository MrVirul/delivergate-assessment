# PROGRESS

## Current status
Phase: 4 — Customers   |   Last gate: PASS (2026-10-27)
Next action: Implement Phase 5 - Seeds (idempotent seeding for users/items/customers)

## Phases
- [x] Phase 0 — Bootstrap
- [x] Phase 1 — Core infrastructure
- [x] Phase 2 — Authentication
- [x] Phase 3 — Inventory
  - [x] DTOs, repo, service, handlers, routes with permissions
- [x] Phase 4 — Customers
  - [x] DTOs, repo, service, handlers, routes with permissions
- [ ] Phase 5 — Seeds
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
