# UI Documentation

## Architecture
- Static files in `web/`, served by Go Fiber
- htmx 2.0.11 vendored locally
- Adapter pattern: `api-ext.js` sends JSON, `render.js` builds HTML
- All dynamic values escaped via `UI.esc`

## Adding a New Screen
Follow the pattern in renderers and add routes/elements.

## Development
Run `make run` and open http://localhost:8080/

## CRUD verification

| Check | Owner | Manager | Cashier |
|---|---|---|---|
| Create item valid | PASS | PASS | FAIL-fixed (403 as per RBAC) |
| Create item invalid (422) | PASS | PASS | N/A |
| Create duplicate item (409) | PASS | PASS | N/A |
| List/search items | PASS | PASS | PASS |
| Edit item PUT/PATCH | PASS | PASS | PASS |
| Delete item | PASS | FAIL-fixed (403 as per RBAC) | FAIL-fixed (403 as per RBAC) |
| Create customer valid | PASS | PASS | PASS |
| Edit/delete customer | PASS (delete), PASS (edit) | PASS (edit only), FAIL-fixed (delete 403) | PASS (edit/create), FAIL-fixed (delete 403) |
| Price 19.99 stored as 1999 cents / shown correctly | PASS | PASS | PASS |
| XSS string renders as plain text | PASS | PASS | PASS |
| 403 toast on forbidden call | PASS | PASS | PASS |
| Expired/removed token redirects to login | PASS | PASS | PASS |
| Server down shows "Cannot reach the server" | PASS | PASS | PASS |
