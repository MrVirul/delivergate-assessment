# PROGRESS-UI

## Current status
Phase: 0 — Wire-up   |   Last gate: PASS (2026-10-27)
Next action: Phase 1 - Create shell, styles, login (index.html, login.html, css, js files)

## Phases
- [x] Phase 0 — Wire-up
  - [x] web/ skeleton
  - [x] htmx 2.0.11 vendored
  - [x] static serving added to router (last)
  - [x] Dockerfile updated
  - [x] tests added
- [ ] Phase 1 — Shell, styles, login
- [ ] Phase 2 — Inventory list
- [ ] Phase 3 — Inventory forms and delete
- [ ] Phase 4 — Customers
- [ ] Phase 5 — Polish and accessibility
- [ ] Phase 6 — Audit and docs

## Blockers
None

## How to run
- go build ./...
- go test -race -count=1 ./...
- make run (or go run ./cmd/api)
- open http://localhost:8080/
