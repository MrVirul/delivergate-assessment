# FRONTEND-AGENT.md — Ice Cream Shop UI (htmx 2.x + CSS + vanilla JavaScript)

> **Audience:** autonomous coding agents (opencode, Antigravity, or any other).
> **Prerequisite:** the Go backend described in `agent.md` is finished and passing its tests.
> **Read first, in this order:** `skills/htmx/SKILL.md` (mandatory), then this file, then `PROGRESS-UI.md` if it exists.
> **Scope:** a simple, fast, good-looking web UI. No framework, no npm build step, no bundler.

---

## 0. How to use this file (for the human)

1. Put this file in the repo root as `frontend-agent.md` and the skill at `skills/htmx/SKILL.md`. Optionally copy the skill to your tool's own skills folder as well (for opencode, `.opencode/skill/htmx/SKILL.md`); the prompts below tell the agent to read it explicitly anyway.
2. Paste the kickoff prompt from section 17.
3. When one tool hits its limit, switch tools and paste the resume prompt. All state lives in `PROGRESS-UI.md`.

---

## 1. Rules of engagement

1. **Do not ask questions.** Decide, then log the decision in `docs/DECISIONS.md` (append; do not overwrite).
2. **The htmx skill is law.** Use **htmx 2.x pinned (2.0.11 or the newest 2.0.x)**, vendored locally. Never htmx 4 and never an unpinned CDN link.
3. **No build tools:** only `.html`, `.css`, `.js` files served as static files.
4. **Do not change backend behavior.** Allowed backend changes: (a) serve the `web/` folder, (b) small CSP/helmet adjustments if the browser console shows violations, (c) tiny DTO tweaks only if the UI truly cannot work otherwise, always with a test. **Never** change the RBAC matrix, routes under `/api`, or response shapes.
5. **Security:** escape every dynamic value (`esc()`), no inline scripts, no inline event handlers, no `hx-on:*`, `allowEval:false`.
6. **No placeholders.** No lorem ipsum, no dead buttons, no `TODO`.
7. **Verify after every phase** (section 12) and commit. Update `PROGRESS-UI.md` before each commit.
8. Stuck >3 attempts on one problem: log it under "Blockers" in `PROGRESS-UI.md`, choose a workaround, continue.
9. Keep files small: one concern per file.

---

## 2. What we are building

A single-page-feeling web app for shop staff, used on a laptop or tablet at the counter and on a phone.

**Screens**
1. **Login** (`/login.html`)
2. **Inventory** (tab) — searchable, filterable, sortable, paginated list; create, edit, delete
3. **Customers** (tab) — searchable, paginated list; create, edit, delete
4. Shell: top bar with shop name, signed-in user + role badge, logout; tab navigation; toast messages; one reusable dialog for forms

**Out of scope:** user management, sales/orders, charts, offline mode, i18n.

---

## 3. Architecture

```
web/
├── index.html            # app shell (requires login)
├── login.html            # login page
├── css/
│   └── styles.css        # all styles, CSS variables, responsive
├── js/
│   ├── render.js         # UI.esc, UI.money, UI.can, UI.coerceForm, UI.render (HTML builders)
│   ├── api-ext.js        # htmx extension "api" (auth header, JSON in/out)
│   └── app.js            # session guard, toasts, dialog glue, global htmx listeners
└── vendor/
    └── htmx.min.js       # pinned 2.x, vendored
```

- **Served by the Go app itself** (same origin, so no CORS): after all `/api` and `/healthz` routes are registered, add Fiber static serving for `./web` with `index.html` as the index. Register it **last** so it never shadows an API route. Update the `Dockerfile` to copy `web/` into the image and make sure the working directory resolves `./web` (or make the path configurable with `WEB_DIR`, default `./web`).
- Add a backend test that proves `GET /` returns 200 HTML, `GET /api/items` without a token still returns 401, and `GET /healthz` still returns 200.
- **Script order** in every page: `vendor/htmx.min.js`, `js/render.js`, `js/api-ext.js`, `js/app.js` (classic scripts, `defer` is fine if the order is kept).
- **htmx config** goes in a `<meta name="htmx-config">` tag exactly as in the skill (section 6): `allowEval:false`, `selfRequestsOnly:true`, `responseHandling` with 204 swappable and 4xx/5xx flagged as errors.
- `<body hx-ext="api">` on both pages.

### Why an adapter?
The API returns JSON, while htmx swaps HTML. `api-ext.js` sends JSON and `UI.render(key, body)` turns JSON into escaped HTML strings. Follow the pattern in `skills/htmx/SKILL.md` section 8. **Test it in the browser's network tab; the reference code is a starting point, not gospel.**

### CSP warning
The backend uses a security-headers middleware (helmet). Its default Content-Security-Policy may block inline scripts, `eval`, or external fonts. After building the shell, open the browser console and fix any "violates Content Security Policy" message by adjusting the server policy narrowly (for example allow `fonts.googleapis.com` and `fonts.gstatic.com` if you use Google Fonts). Do not disable the middleware.

---

## 4. Session and permissions in the browser

- On successful login the adapter stores in `sessionStorage`: `token`, `user` (JSON: id, name, email, role), `expires_at` (now + `expires_in`).
- **Guard:** `app.js` on `index.html` runs first thing: if there is no token or it is expired, `location.replace('/login.html')`. On `login.html`, if a valid token exists, go to `/`.
- **401 anywhere** (`htmx:responseError` with status 401): clear the session, redirect to login with a toast message "Your session has ended. Please sign in again."
- **Logout button:** `hx-post="/api/auth/logout"` with `hx-swap="none"`; in an `htmx:afterRequest` listener, clear the session and redirect **whether or not the request succeeded**.
- **UI permissions (UX only, the server enforces for real):** `UI.can(permission)` mirrors the matrix below. Hide controls the role cannot use. If a forbidden call still happens (403), show a toast "You don't have permission to do that."

| Control | Owner | Manager | Cashier |
|---|:-:|:-:|:-:|
| See inventory list and details | ✅ | ✅ | ✅ |
| "Add item" button | ✅ | ✅ | hidden |
| "Edit" on an item | ✅ | ✅ | ✅ |
| "Delete" on an item | ✅ | hidden | hidden |
| See customers | ✅ | ✅ | ✅ |
| "Add customer", "Edit" | ✅ | ✅ | ✅ |
| "Delete" on a customer | ✅ | hidden | hidden |

---

## 5. API reference the UI uses

Base: same origin. Envelope: success `{ "data": ..., "meta": {...} }`; error `{ "error": { "code", "message", "details": [{ "field", "message" }] } }`. Money is **integer cents** (`price_cents`).

| UI action | Request | `data-render` key |
|---|---|---|
| Login | `POST /api/auth/login` `{email,password}` | `login` (stores session, returns empty string) |
| Logout | `POST /api/auth/logout` (204) | none |
| Items list | `GET /api/items?search=&category=&low_stock=&sort=&page=&limit=` | `items-list` |
| Open "add item" form | none (static form from `UI.render('item-form', null)` via a local template call, or `hx-get` to a static partial) | `item-form` |
| Open "edit item" form | `GET /api/items/:id` | `item-form` |
| Save new item | `POST /api/items` | none (`hx-swap="none"`) |
| Save edited item | `PUT /api/items/:id` (all fields) | none |
| Delete item | `DELETE /api/items/:id` (204) | none; row removed with `hx-target="closest tr" hx-swap="outerHTML"` or list refresh |
| Customers list | `GET /api/customers?search=&sort=&page=&limit=` | `customers-list` |
| Customer forms and save/delete | same pattern as items | `customer-form` |

Item fields: `name`, `category` (`flavor`, `topping`, `cone`, `cup`, `packaging`, `other`), `description`, `quantity`, `unit` (`tub`, `scoop`, `kg`, `litre`, `pcs`, `box`), `price_cents`, `reorder_level`, computed `low_stock`.
Customer fields: `name`, `email` (optional), `phone`, `address` (optional), `notes` (optional).
Sort values: items `name`, `quantity`, `price_cents`, `created_at` (prefix `-` for descending); customers `name`, `created_at`.

**Read `docs/openapi.yaml` and the DTO validation tags before writing forms** so field names, limits and required flags match exactly. If something differs from this table, the backend code wins; log it in `docs/DECISIONS.md`.

---

## 6. Behavior details

### Lists
- Container pattern: `<div id="items-results" hx-get="/api/items" hx-include="#items-filters" hx-trigger="load, items-changed from:body" data-render="items-list">`.
- Filters form (`#items-filters`): text search (debounced 300 ms, `hx-sync="this:replace"`), category select, "Low stock only" checkbox, sort select. Any change reloads the list from page 1.
- Pagination: renderer outputs Previous/Next buttons and "Showing 21–40 of 57" from `meta`. Buttons use `hx-vals='{"page": N}'` plus `hx-include` of the filters form. Disable Previous on page 1 and Next on the last page.
- Empty state: if there are no results, say so and offer the next step ("No items match your search. Clear filters." / "No customers yet. Add the first one.").
- Loading: show a subtle indicator during requests (`hx-indicator`, `htmx-request` class). No layout jump.

### Forms (dialog)
- One `<dialog id="modal">` in `index.html` with a body container `#modal-body`.
- "Add" and "Edit" buttons load the form HTML into `#modal-body` (`hx-target="#modal-body"`, `hx-swap="innerHTML"`); `app.js` opens the dialog in an `htmx:afterSwap` listener when the target is `#modal-body`, and focuses the first field.
- Form submit: `hx-post` (create) or `hx-put` (edit), `hx-swap="none"`, `hx-disabled-elt="find button[type=submit]"`.
- On `htmx:afterRequest` with `successful`: close the dialog, show a success toast ("Item saved" / "Customer saved"), and `htmx.trigger(document.body, 'items-changed')` or `'customers-changed'`.
- On 422: show the server's `details` next to the matching fields (`aria-invalid="true"`, message in a `<p class="field-error">` linked with `aria-describedby`) and keep the dialog open. On 409: show the message at the top of the form. On other errors: toast.
- Price input is a decimal in the main currency (for example `185.00`); `data-money` makes `UI.coerceForm` convert it to cents. When editing, show `price_cents / 100` with two decimals.
- Escape closes the dialog; clicking the backdrop closes it; focus returns to the button that opened it.

### Delete
- Button has `hx-delete`, `hx-confirm="Delete “Chocolate”? This cannot be undone."` (the name must be escaped in the attribute). On success remove the row (or reload the list) and toast "Item deleted". On failure toast the message.

### Currency
`UI.money(cents)` uses `Intl.NumberFormat` with a single constant `CURRENCY = 'LKR'` at the top of `render.js` (easy to change). Dates via `Intl.DateTimeFormat`.

### Toasts
`UI.toast(message, kind)` with kinds `success`, `error`, `info`. Container has `role="status" aria-live="polite"`. Auto-dismiss after 4 s; errors after 6 s; each has a close button.

---

## 7. Design direction (decided; follow it)

Subject: a neighbourhood ice cream shop's back office. Audience: busy staff, often on a tablet with sticky fingers, at a counter. Primary job: find an item or a customer fast, and change stock without mistakes.

**Palette (CSS variables on `:root`)**

| Token | Hex | Use |
|---|---|---|
| `--freezer` | `#F2F6F8` | page background (cool, like a freezer case) |
| `--panel` | `#FFFFFF` | cards, table, dialog |
| `--cocoa` | `#2E2321` | text |
| `--cocoa-soft` | `#6B5B57` | secondary text |
| `--berry` | `#C2185B` | primary actions, active tab |
| `--mint` | `#1E9E78` | healthy stock, success |
| `--butterscotch` | `#E09A1A` | low stock, warnings |
| `--line` | `#D9E2E7` | borders |

Also define a darker `--berry-ink` for hover/pressed, and make sure all text/background pairs meet WCAG AA contrast.

**Type:** one display face for headings and the shop name (`"Bricolage Grotesque"` from Google Fonts, or a local system fallback if you prefer no network), and `system-ui` for everything else. Line length under 80 characters, sentence case, no all-caps labels, no tracked-out eyebrow text.

**The one memorable element:** a **stock gauge** on every inventory row: a slim horizontal bar showing quantity against the reorder level (`fill % = min(100, quantity / (reorder_level × 3 or 20) × 100)`), mint when stock is healthy, butterscotch with the text "Low stock" when `low_stock` is true. Keep everything else quiet and consistent. Do not add gradients, glassmorphism, or decorative animation.

**Layout**
```
┌──────────────────────────────────────────────────────┐
│ Scoop & Co.                     Nimal · Manager  [Log out] │
├──────────────────────────────────────────────────────┤
│ [ Inventory ] [ Customers ]                          │
├──────────────────────────────────────────────────────┤
│ [search…………] [Category ▾] [☐ Low stock] [Sort ▾]  [Add item] │
│ ┌──────────────────────────────────────────────────┐ │
│ │ Name        Category  Stock gauge   Price  Actions │ │
│ │ Chocolate   flavor    ███████░ 12   1,850  Edit Delete │ │
│ └──────────────────────────────────────────────────┘ │
│ Showing 1–20 of 57                  [Previous] [Next] │
└──────────────────────────────────────────────────────┘
```
- Desktop/tablet: table. **Phone (<720 px): each row becomes a stacked card** (CSS only: `display:block` rows with `data-label` attributes, no duplicate markup).
- Touch targets at least 44×44 px. Visible `:focus-visible` outlines. Respect `prefers-reduced-motion`. Left-align text; right-align numbers.
- Use the shop name "Scoop & Co." as a placeholder constant in one place (`<title>`, top bar) so it is easy to change.
- Write button labels as specific actions: "Add item", "Save item", "Save customer", "Delete", "Sign in". Error text says what happened and what to do.

---

## 8. Phases and acceptance criteria

Track tasks in `PROGRESS-UI.md` (format in section 14). Times are a guide.

### Phase 0 — Wire-up (≈ 45 min)
- Create `web/` skeleton; vendor htmx 2.x (see skill section 0); add static serving to the Go app; update Dockerfile and README note; add the backend static/API-priority test.
- **Accept:** `make run` then `curl -i localhost:8080/` returns 200 HTML; `/api/items` still 401 without a token; `go test -race ./...` green; `htmx.version` shows `2.x` in the console.

### Phase 1 — Shell, styles, login (≈ 1.5 h)
- `styles.css` with tokens, base typography, layout, buttons, form controls, toast, dialog, table/card responsive rules.
- `login.html` with htmx login; `app.js` guard, toast, logout; `index.html` shell with tabs (tab switching with plain `hx-get` of the list into the main panel or simple show/hide; choose the simplest that keeps the URL hash `#inventory` / `#customers` working on reload).
- **Accept:** each seeded user can sign in and out; wrong password shows a clear error; visiting `/` without a session redirects to login; the role badge shows the right role.

### Phase 2 — Inventory list (≈ 2 h)
- `items-list` renderer with stock gauge, filters, sort, pagination, empty/loading states.
- **Accept:** search, category filter, low-stock filter, sort and paging all work against the real API; no XSS (create an item named `<img src=x onerror=alert(1)>` and confirm it renders as text).

### Phase 3 — Inventory forms and delete (≈ 1.5 h)
- Dialog create/edit, validation display, delete with confirm, role-based buttons.
- **Accept:** Owner can add/edit/delete; Manager sees no Delete; Cashier sees no Add or Delete but can Edit; forcing a forbidden call (via DevTools) shows the 403 toast.

### Phase 4 — Customers (≈ 1.5 h)
- Reuse the same structure for customers (list, search, forms, delete).
- **Accept:** same criteria as phases 2–3 for the customer matrix; phone and email validation errors appear next to the fields.

### Phase 5 — Polish and accessibility (≈ 1.5 h)
- Keyboard-only walkthrough, focus return after dialog, screen-reader labels, reduced motion, contrast check, 360 px phone layout, 768 px tablet layout, long names do not break layout.
- Handle server down (`htmx:sendError`) and expired tokens gracefully.

### Phase 6 — Audit and docs (≈ 1 h)
- Run section 13. Add a "Web UI" section to `README.md` (how to open it, seeded logins, folder map) and a short `docs/UI.md` (architecture, adapter pattern, how to add a new screen). Take screenshots if the environment allows and save in `docs/screenshots/`.

---

## 9. Reference snippets (follow these shapes)

**`render.js` skeleton**
```js
const CURRENCY = 'LKR';
const PERMISSIONS = {
  Owner:   ['items:read','items:create','items:update','items:delete','customers:read','customers:create','customers:update','customers:delete'],
  Manager: ['items:read','items:create','items:update','customers:read','customers:create','customers:update'],
  Cashier: ['items:read','items:update','customers:read','customers:create','customers:update'],
};

const UI = {
  esc: (v) => String(v ?? '').replace(/[&<>"']/g, (c) =>
    ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c])),
  money: (cents) => new Intl.NumberFormat(undefined, { style: 'currency', currency: CURRENCY }).format(cents / 100),
  user: () => JSON.parse(sessionStorage.getItem('user') || 'null'),
  can: (perm) => (PERMISSIONS[UI.user()?.role] || []).includes(perm),
  // coerceForm(formData, form) -> plain object (numbers, cents, nulls) — see skills/htmx/SKILL.md section 8
  // render(key, body) -> HTML string using RENDERERS[key]
};
const RENDERERS = { /* 'login', 'items-list', 'item-form', 'customers-list', 'customer-form' */ };
```

**Item row (shape, adapt names)**
```js
const row = (it) => `
  <tr>
    <td data-label="Name">${UI.esc(it.name)}</td>
    <td data-label="Category">${UI.esc(it.category)}</td>
    <td data-label="Stock">${gauge(it)}</td>
    <td data-label="Price">${UI.money(it.price_cents)}</td>
    <td class="actions">
      ${UI.can('items:update') ? editButton(it) : ''}
      ${UI.can('items:delete') ? deleteButton(it) : ''}
    </td>
  </tr>`;
```
For attributes that contain user text (for example `hx-confirm`), escape with `UI.esc` as well.

---

## 10. Testing

Automated (Go, in the existing suite):
- `GET /` serves HTML; `/api/*` and `/healthz` unaffected by static serving.

Manual QA matrix (record results in `docs/UI.md`):

| Check | Owner | Manager | Cashier |
|---|---|---|---|
| Sign in with seeded account | ☐ | ☐ | ☐ |
| Add item button visible | ☐ yes | ☐ yes | ☐ hidden |
| Edit item works | ☐ | ☐ | ☐ |
| Delete item button visible | ☐ yes | ☐ hidden | ☐ hidden |
| Add / edit customer works | ☐ | ☐ | ☐ |
| Delete customer button visible | ☐ yes | ☐ hidden | ☐ hidden |
| Direct forbidden call shows 403 toast | n/a | ☐ | ☐ |

Plus: expired/removed token redirects to login; server stopped shows "Cannot reach the server"; XSS test string renders as text; 360 px layout usable; keyboard-only flow works; console has no errors or CSP violations.

If Playwright or a similar tool is already available you may add a short smoke script, but it is optional.

---

## 11. Git workflow

- Conventional commits: `feat(ui): items list with stock gauge`, `fix(ui): ...`, `docs: ...`.
- One commit per completed task group; never commit a broken build; update `PROGRESS-UI.md` in the same commit.
- Final: tag `v1.1.0`.

---

## 12. Verification gate (after every phase)

```bash
gofmt -l .
go vet ./...
go build ./...
go test -race -count=1 ./...
make run &           # then:
curl -s -o /dev/null -w "%{http_code}\n" localhost:8080/            # 200
curl -s -o /dev/null -w "%{http_code}\n" localhost:8080/api/items   # 401
curl -s -o /dev/null -w "%{http_code}\n" localhost:8080/healthz     # 200
```
Then open the page in a browser (or headless browser if available), exercise the new feature, and check the console for errors. Do not move on while anything is red.

---

## 13. Final audit

- [ ] htmx file is vendored, pinned 2.x, and `htmx.version` confirms it; no CDN scripts.
- [ ] No inline scripts, no inline `on*` handlers, no `hx-on:*`; `allowEval:false` set.
- [ ] Every dynamic value in HTML strings goes through `UI.esc` (grep `render.js` for `${` and review each one).
- [ ] All three roles verified in the browser against the matrix in section 4.
- [ ] 401 handling, 403 toast, 409 and 422 messages, network error all verified.
- [ ] Works at 360 px, 768 px and desktop widths; keyboard-only flow works; focus returns after the dialog.
- [ ] Backend tests still green; static serving does not shadow any API route.
- [ ] README and `docs/UI.md` updated; QA matrix filled in.
- [ ] `grep -rn "TODO\|FIXME\|console.log" web/` is empty.

---

## 14. `PROGRESS-UI.md` format

```markdown
# PROGRESS-UI

## Current status
Phase: <N> — <name>   |   Last gate: PASS/FAIL (<date time>)
Next action: <one concrete sentence an agent can start on immediately>

## Phases
- [ ] Phase 0 — Wire-up
- [ ] Phase 1 — Shell, styles, login
...

## Blockers
- <date> <what failed> <what was tried> <workaround>

## How to run
<exact commands that work right now>
```
Rules: update before every commit; tick only what passed the gate; "Next action" must always be accurate.

---

## 15. Definition of done

1. Section 13 audit fully checked.
2. A clean clone, `make migrate && make seed && make run`, then opening `http://localhost:8080` lets all three seeded users sign in and do exactly what their role allows.
3. Backend and (optional) UI tests green; README and `docs/UI.md` complete; repo tagged `v1.1.0`.

---

## 16. Common traps (read before coding)

- Number inputs arrive as strings: convert in `UI.coerceForm` or the API returns 400/422.
- 204 responses do not swap unless `responseHandling` allows it.
- 4xx/5xx are ignored unless flagged `error:true`.
- The security middleware's CSP can silently break inline code and fonts: watch the console.
- Do not use htmx 4 event names or docs.
- Static serving registered before the API routes would shadow them: register it last.
- Hiding a button is not security; the server decides, the UI only reflects it.

---

## 17. Prompts for the human to paste

**Kickoff (opencode):**
> Read `skills/htmx/SKILL.md` completely, then read `frontend-agent.md` completely. The Go backend from `agent.md` is already finished. Create `PROGRESS-UI.md` as specified in section 14 and start Phase 0. Work autonomously through the phases in order. After each phase run the verification gate (section 12), update `PROGRESS-UI.md`, and commit. Do not ask me questions; make the simplest reasonable decision and log it in `docs/DECISIONS.md`. Use htmx 2.x only. Do not change backend behavior beyond what `frontend-agent.md` allows. Before stopping for any reason, make sure the "Next action" line is accurate and everything is committed.

**Resume (new session or Antigravity):**
> Read `skills/htmx/SKILL.md`, `frontend-agent.md`, then `PROGRESS-UI.md` and `git log --oneline -20`. Run the verification gate first and fix anything red. Then continue from the "Next action" line. Same rules: no questions, log decisions, update `PROGRESS-UI.md` before every commit.

**Fix prompt:**
> Something is broken: <describe, and paste the console or network output>. Find the root cause (check DevTools Network first, then the Console for CSP or JS errors), fix it, and verify in the browser. Do not weaken security settings to make it work.

**Audit prompt:**
> Execute section 13 of `frontend-agent.md` item by item. For each item show the command or browser check you used and the result. Fix every failure and re-verify.
