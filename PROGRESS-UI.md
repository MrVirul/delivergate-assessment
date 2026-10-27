# PROGRESS-UI

## Current status
Phase: 1 — Shell, styles, login   |   Last gate: PASS (2026-10-27)
Next action: Phase 2 - Inventory list (ensure renderers work well, test with API)

## Phases
- [x] Phase 0 — Wire-up
- [x] Phase 1 — Shell, styles, login
  - [x] styles.css with tokens
  - [x] login.html with htmx login
  - [x] index.html shell with tabs
  - [x] app.js guard, toasts, logout
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
- make run
- open http://localhost:8080/
