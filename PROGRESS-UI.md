# PROGRESS-UI

## Current status
Phase: 6 — Audit and docs   |   Last gate: PASS (2026-10-27)
Next action: All UI phases complete

## Phases
- [x] Phase 0 — Wire-up
- [x] Phase 1 — Shell, styles, login
- [x] Phase 2 — Inventory list
- [x] Phase 3 — Inventory forms and delete
- [x] Phase 4 — Customers
- [x] Phase 5 — Polish and accessibility
- [x] Phase 6 — Audit and docs
  - [x] docs/UI.md, README updated
  - [x] gates pass

## Blockers
None

## How to run
- go build ./...
- go test -race -count=1 ./...
- make run
- open http://localhost:8080/

## Recent updates
- Implemented custom confirmation card (js/confirm.js) loaded after app.js with htmx:confirm handling, dialog stacking, keyboard/accessibility rules
- Switched delete buttons to use htmx delete with confirmation card
- Added edit confirmation with change detection; shows "No changes to save." if no changes
- Added login.js separation and adjusted login flow
- Updated UI docs with CRUD verification table
