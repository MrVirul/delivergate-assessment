# PROGRESS-UI

## Current status
Phase: 4 — Customers   |   Last gate: PASS (2026-10-27)
Next action: Phase 5 - Polish and accessibility

## Phases
- [x] Phase 0 — Wire-up
- [x] Phase 1 — Shell, styles, login
- [x] Phase 2 — Inventory list
- [x] Phase 3 — Inventory forms and delete
- [x] Phase 4 — Customers
- [ ] Phase 5 — Polish and accessibility
- [ ] Phase 6 — Audit and docs

## Blockers
None

## How to run
- go build ./...
- go test -race -count=1 ./...
- make run
- open http://localhost:8080/
