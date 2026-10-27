# PROGRESS

## Current status
Phase: 2 — Authentication   |   Last gate: PASS (2026-10-27)
Next action: Implement Phase 3 - Inventory CRUD (models/repos/services/handlers/routes with permissions)

## Phases
- [x] Phase 0 — Bootstrap
- [x] Phase 1 — Core infrastructure
- [x] Phase 2 — Authentication
  - [x] auth (permissions, jwt, password)
  - [x] repos (users, revoked tokens)
  - [x] service (login/logout)
  - [x] middleware (auth, rbac, errors)
  - [x] handlers/routes (login, logout, healthz)
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
- make migrate
- make run
- docker compose up -d db
