# AGENT.md — Ice Cream Shop Inventory Management System (Go + Fiber + GORM)

> **Audience:** autonomous coding agents (opencode, Antigravity, or any other). You are the only developer.
> **Deadline:** 2 days total. Work fast, but never ship something unverified.
> **Read this whole file first. Then read `PROGRESS.md` (if it exists) and continue from where it stops.**

---

## 0. How to use this file (for the human)

1. Put this file in an empty repo root as `agent.md` (also copy it to `AGENTS.md` — opencode auto-loads that name; if your tool does not auto-load it, start with the kickoff prompt in section 20).
2. Paste the kickoff prompt (section 20). The agent works phase by phase (section 12), ticking tasks in `PROGRESS.md`.
3. When opencode hits a limit, open Antigravity in the same repo and paste the **resume prompt** (section 20). All state lives in files, never in chat memory.

---

## 1. Your role and operating rules

You are a senior Go backend engineer (20 years of experience). You write production-quality, boring, readable code.

**Hard rules**

1. **Do not ask the human questions.** Every ambiguity is already resolved in this file. If something is truly undefined, make the simplest reasonable decision, and append it to `docs/DECISIONS.md` (date, decision, reason).
2. **Spec is law.** The RBAC matrix (section 4) must be implemented exactly. Do not "improve" permissions.
3. **No placeholders.** No `TODO`, no `panic("not implemented")`, no fake data in handlers, no commented-out code in the final result.
4. **Verify after every phase** with the Verification Gate (section 13). Never start the next phase while the gate is red.
5. **Stay inside the approved dependency list** (section 3). If you need something else, prefer the standard library.
6. **Small, frequent commits** (section 17). Update `PROGRESS.md` before every commit.
7. **Never hardcode secrets.** All config comes from environment variables. `.env` is git-ignored; `.env.example` is committed.
8. **Never invent APIs from memory.** If unsure about a library function signature, read the module source with `go doc <pkg> <symbol>` or look in the module cache (`$(go env GOMODCACHE)`), then write the code.
9. **If stuck >3 attempts on the same error:** write the blocker in `PROGRESS.md` under "Blockers", choose a pragmatic workaround, and move on. Do not loop.
10. **Context hygiene:** keep files small (<300 lines where possible), keep functions focused, and keep `PROGRESS.md` current so another agent can resume with zero chat history.

---

## 2. Objective (from the assignment)

Backend for an ice cream shop managing **inventory items** and **customer records**, with **JWT authentication** and **role-based access control**. Three roles: `Owner`, `Manager`, `Cashier`.

**Deliverables (all mandatory):**

1. Functional backend application
2. Authentication + role-based authorization
3. Inventory CRUD APIs
4. Customer CRUD APIs
5. Database schema + ER diagram
6. Migrations (or equivalent setup mechanism)
7. Seeded initial accounts for the three roles
8. Unit tests covering core functionality **and access permissions**
9. API documentation (endpoints, request/response formats, authorization requirements)

---

## 3. Technology decisions (final — do not revisit)

| Concern | Choice | Notes |
|---|---|---|
| Language | Go 1.22+ | Use `go.mod` with the installed toolchain version |
| Web framework | `github.com/gofiber/fiber/v2` | Use v2 (stable, widely documented). Not v3 |
| ORM | `gorm.io/gorm` | Required by spec |
| Prod/dev DB | PostgreSQL 16 via `gorm.io/driver/postgres` | Run with docker-compose |
| Test DB | In-memory SQLite via `github.com/glebarez/sqlite` | Pure Go driver, **no CGO**, works on any machine |
| JWT | `github.com/golang-jwt/jwt/v5` | HS256 |
| Password hashing | `golang.org/x/crypto/bcrypt` | cost 12 (cost 4 in tests for speed) |
| Validation | `github.com/go-playground/validator/v10` | struct tags |
| Env loading | `github.com/joho/godotenv` | optional `.env` |
| Test helpers | `github.com/stretchr/testify` | `assert`, `require` |
| Logging | stdlib `log/slog` | JSON handler |

Fiber middleware allowed from `github.com/gofiber/fiber/v2/middleware/*`: `recover`, `requestid`, `logger`, `cors`, `helmet`, `limiter`.

**Migrations decision:** GORM `AutoMigrate` in a dedicated `internal/database/migrate.go`, executed via `cmd/migrate` and (optionally) on server start when `AUTO_MIGRATE=true`. It works on both Postgres and SQLite, which keeps tests honest. The ER diagram in `docs/ER_DIAGRAM.md` is the design source of truth.

**Identifier decision:** `uint` auto-increment primary keys (clean `{id}` paths).
**Money decision:** store prices as integer **cents** (`price_cents int64`) — never floats.
**Deletion decision:** hard delete (the spec says "delete"; keeps unique constraints simple).

---

## 4. Authorization matrix (single source of truth)

| Operation | Owner | Manager | Cashier |
|---|:-:|:-:|:-:|
| View inventory items | ✅ | ✅ | ✅ |
| Add inventory items | ✅ | ✅ | ❌ |
| Edit inventory items | ✅ | ✅ | ✅ |
| Delete inventory items | ✅ | ❌ | ❌ |
| View customer records | ✅ | ✅ | ✅ |
| Add customer records | ✅ | ✅ | ✅ |
| Edit customer records | ✅ | ✅ | ✅ |
| Delete customer records | ✅ | ❌ | ❌ |

Implement as permissions, not scattered `if role ==` checks:

```go
// internal/auth/permissions.go
type Role string

const (
    RoleOwner   Role = "Owner"
    RoleManager Role = "Manager"
    RoleCashier Role = "Cashier"
)

type Permission string

const (
    ItemsRead       Permission = "items:read"
    ItemsCreate     Permission = "items:create"
    ItemsUpdate     Permission = "items:update"
    ItemsDelete     Permission = "items:delete"
    CustomersRead   Permission = "customers:read"
    CustomersCreate Permission = "customers:create"
    CustomersUpdate Permission = "customers:update"
    CustomersDelete Permission = "customers:delete"
)

var rolePermissions = map[Role]map[Permission]struct{}{
    RoleOwner: set(ItemsRead, ItemsCreate, ItemsUpdate, ItemsDelete,
        CustomersRead, CustomersCreate, CustomersUpdate, CustomersDelete),
    RoleManager: set(ItemsRead, ItemsCreate, ItemsUpdate,
        CustomersRead, CustomersCreate, CustomersUpdate),
    RoleCashier: set(ItemsRead, ItemsUpdate,
        CustomersRead, CustomersCreate, CustomersUpdate),
}

func Can(role Role, p Permission) bool { _, ok := rolePermissions[role][p]; return ok }
```

Unknown or empty role → deny everything.

---

## 5. Repository layout

```
icecream-inventory/
├── agent.md / AGENTS.md
├── PROGRESS.md                  # live task tracker (agents update this)
├── README.md
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── .gitignore
├── .golangci.yml                # optional, small config
├── go.mod / go.sum
├── cmd/
│   ├── api/main.go              # server entrypoint
│   ├── migrate/main.go          # runs migrations
│   └── seed/main.go             # runs seeders
├── internal/
│   ├── config/config.go         # env -> Config struct, validation
│   ├── database/
│   │   ├── database.go          # Open(cfg) for postgres; OpenSQLite for tests
│   │   ├── migrate.go           # AutoMigrate(models...)
│   │   └── seed.go              # idempotent seeders
│   ├── models/                  # GORM models
│   │   ├── user.go
│   │   ├── item.go
│   │   ├── customer.go
│   │   └── revoked_token.go
│   ├── auth/
│   │   ├── permissions.go
│   │   ├── jwt.go               # token issue/parse
│   │   └── password.go          # hash/compare
│   ├── dto/                     # request/response structs + validation tags
│   ├── repository/              # DB access (interfaces + gorm impl)
│   ├── service/                 # business logic
│   ├── handler/                 # Fiber handlers (thin)
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── rbac.go              # RequirePermission(...)
│   │   └── errors.go            # central error handler
│   ├── httpx/                   # response envelope, error types, pagination, validator
│   ├── router/router.go         # route wiring
│   └── testutil/                # test app builder, tokens, seeded users
├── docs/
│   ├── API.md                   # human API documentation
│   ├── openapi.yaml             # OpenAPI 3.0 spec
│   ├── ER_DIAGRAM.md            # Mermaid ER diagram + entity descriptions
│   ├── DECISIONS.md             # decision log
│   └── TESTING.md               # how to run tests + coverage matrix
└── .github/workflows/ci.yml     # optional: build, vet, test
```

Layering rule: **handler → service → repository → DB**. Handlers never touch GORM. Services never import Fiber. Repositories are defined as interfaces in `repository` so services are unit-testable with fakes.

---

## 6. Configuration

`.env.example` (commit this; never commit `.env`):

```
APP_ENV=development            # development | production | test
APP_PORT=8080
DATABASE_URL=postgres://icecream:icecream@localhost:5432/icecream?sslmode=disable
JWT_SECRET=change-me-to-a-random-string-of-at-least-32-characters
JWT_TTL_MINUTES=60
AUTO_MIGRATE=true
BCRYPT_COST=12
CORS_ALLOW_ORIGINS=http://localhost:3000
LOGIN_RATE_LIMIT_PER_MIN=10

# Seed account passwords (dev defaults; override in any shared environment)
SEED_OWNER_EMAIL=owner@icecream.local
SEED_OWNER_PASSWORD=Owner@12345
SEED_MANAGER_EMAIL=manager@icecream.local
SEED_MANAGER_PASSWORD=Manager@12345
SEED_CASHIER_EMAIL=cashier@icecream.local
SEED_CASHIER_PASSWORD=Cashier@12345
```

`config.Load()` must: read env (load `.env` if present), apply defaults, and **fail fast** with a clear error if `JWT_SECRET` is shorter than 32 chars (except `APP_ENV=test`) or `DATABASE_URL` is missing outside test.

---

## 7. Data model

### 7.1 Entities

**users**
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| name | varchar(100) | not null |
| email | varchar(255) | not null, **unique**, stored lowercase |
| password_hash | varchar(255) | not null, never serialized (`json:"-"`) |
| role | varchar(20) | not null, one of Owner/Manager/Cashier (check in code + DB check if portable) |
| created_at, updated_at | timestamp | |

**inventory_items**
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| name | varchar(120) | not null, **unique** (case-insensitive at service level) |
| category | varchar(30) | not null, one of `flavor`, `topping`, `cone`, `cup`, `packaging`, `other` |
| description | text | nullable |
| quantity | int | not null, `>= 0` |
| unit | varchar(20) | not null, one of `tub`, `scoop`, `kg`, `litre`, `pcs`, `box` |
| price_cents | bigint | not null, `>= 0` |
| reorder_level | int | not null, default 0, `>= 0` |
| created_by_id | uint | FK → users.id, nullable, ON DELETE SET NULL |
| updated_by_id | uint | FK → users.id, nullable, ON DELETE SET NULL |
| created_at, updated_at | timestamp | |

Computed (not stored): `low_stock = quantity <= reorder_level` (include in the response).
Index: `name` (unique), `category`.

**customers**
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| name | varchar(120) | not null |
| email | varchar(255) | nullable, **unique when present** (use `*string`) |
| phone | varchar(30) | not null |
| address | varchar(255) | nullable |
| notes | text | nullable |
| created_by_id | uint | FK → users.id, nullable, ON DELETE SET NULL |
| updated_by_id | uint | FK → users.id, nullable, ON DELETE SET NULL |
| created_at, updated_at | timestamp | |

Index: `phone`, `name`.

**revoked_tokens** (makes logout real for stateless JWT)
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| jti | varchar(64) | not null, unique |
| expires_at | timestamp | not null, indexed (allows cleanup) |
| created_at | timestamp | |

### 7.2 Relationships

- `users` 1 — N `inventory_items` (via `created_by_id`, `updated_by_id`)
- `users` 1 — N `customers` (via `created_by_id`, `updated_by_id`)
- `revoked_tokens` is standalone (references a token id, not a user FK)

### 7.3 ER diagram deliverable

Create `docs/ER_DIAGRAM.md` with a Mermaid `erDiagram` that matches the models **exactly**, plus a table describing each entity. Example skeleton (extend with all columns):

```mermaid
erDiagram
    USERS ||--o{ INVENTORY_ITEMS : "creates/updates"
    USERS ||--o{ CUSTOMERS : "creates/updates"
    USERS {
        uint id PK
        string name
        string email UK
        string password_hash
        string role
        datetime created_at
        datetime updated_at
    }
    INVENTORY_ITEMS {
        uint id PK
        string name UK
        string category
        text description
        int quantity
        string unit
        bigint price_cents
        int reorder_level
        uint created_by_id FK
        uint updated_by_id FK
    }
    CUSTOMERS {
        uint id PK
        string name
        string email UK
        string phone
        string address
        text notes
        uint created_by_id FK
        uint updated_by_id FK
    }
    REVOKED_TOKENS {
        uint id PK
        string jti UK
        datetime expires_at
    }
```

---

## 8. API specification

Base path: `/api`. All bodies are JSON. Content-Type `application/json`.

### 8.1 Response envelope

Success (single):
```json
{ "data": { ... } }
```
Success (list):
```json
{ "data": [ ... ], "meta": { "page": 1, "limit": 20, "total": 57, "total_pages": 3 } }
```
Error:
```json
{ "error": { "code": "VALIDATION_ERROR", "message": "Validation failed", "details": [ { "field": "quantity", "message": "must be >= 0" } ] } }
```
`details` is only present for validation errors. Never leak internal error text, SQL, or stack traces.

### 8.2 Error codes → HTTP status

| Situation | Status | `code` |
|---|---|---|
| Malformed JSON / bad path id / bad query param | 400 | `BAD_REQUEST` |
| Missing/invalid/expired/revoked token, bad credentials | 401 | `UNAUTHORIZED` |
| Authenticated but role lacks permission | 403 | `FORBIDDEN` |
| Resource not found | 404 | `NOT_FOUND` |
| Unique conflict (item name, customer email, etc.) | 409 | `CONFLICT` |
| Validation failure | 422 | `VALIDATION_ERROR` |
| Rate limit hit | 429 | `RATE_LIMITED` |
| Anything unexpected | 500 | `INTERNAL_ERROR` |

### 8.3 Check order (important — tests rely on it)

For protected routes: **authenticate (401) → authorize by role (403) → parse/validate body (400/422) → load resource (404) → business rules (409)**.
Example: a Cashier calling `DELETE /api/items/99999` gets **403**, not 404.

### 8.4 Endpoints

**Auth**

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/login` | public | Body `{email, password}` → `{data:{access_token, token_type:"Bearer", expires_in, user:{id,name,email,role}}}` |
| POST | `/api/auth/logout` | any authenticated role | Revokes current token's `jti`. Returns **204** |

- Login failure: always **401** with the same message `"invalid email or password"` (no user enumeration). When the email does not exist, still run a bcrypt compare against a dummy hash to equalize timing.
- Login is rate limited per IP (`LOGIN_RATE_LIMIT_PER_MIN`).

**Inventory**

| Method | Path | Permission | Notes |
|---|---|---|---|
| GET | `/api/items` | `items:read` | Query: `search` (name, case-insensitive contains), `category`, `low_stock=true`, `page` (default 1), `limit` (default 20, max 100), `sort` (`name`,`quantity`,`price_cents`,`created_at`; prefix `-` for desc, default `name`) |
| GET | `/api/items/:id` | `items:read` | 404 if missing |
| POST | `/api/items` | `items:create` | 201 + created item |
| PUT | `/api/items/:id` | `items:update` | Full replace; all required fields must be present |
| PATCH | `/api/items/:id` | `items:update` | Partial update; at least one field required |
| DELETE | `/api/items/:id` | `items:delete` | 204 |

**Customers**

| Method | Path | Permission | Notes |
|---|---|---|---|
| GET | `/api/customers` | `customers:read` | Query: `search` (name/email/phone contains), `page`, `limit`, `sort` (`name`,`created_at`) |
| GET | `/api/customers/:id` | `customers:read` | |
| POST | `/api/customers` | `customers:create` | 201 |
| PUT | `/api/customers/:id` | `customers:update` | Full replace |
| PATCH | `/api/customers/:id` | `customers:update` | Partial |
| DELETE | `/api/customers/:id` | `customers:delete` | 204 |

**Extras**

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/healthz` | public | `{ "status": "ok" }` — also pings the DB, 503 if DB down |

### 8.5 Validation rules

Item create/PUT:
- `name`: required, trimmed, 2–120 chars
- `category`: required, one of the allowed values
- `description`: optional, max 1000
- `quantity`: required, integer ≥ 0 (use pointer types in DTO so `0` is distinguishable from missing)
- `unit`: required, one of allowed values
- `price_cents`: required, integer ≥ 0
- `reorder_level`: optional, integer ≥ 0 (default 0)

Customer create/PUT:
- `name`: required, trimmed, 2–120
- `email`: optional, valid email, stored lowercase, unique when present
- `phone`: required, 7–30 chars, digits and `+ - ( ) space` only
- `address`: optional, max 255
- `notes`: optional, max 1000

PATCH DTOs use pointer fields; only non-nil fields are applied. Reject unknown fields? **No** — ignore them, but never let clients set `id`, `created_by_id`, timestamps, or `role`.

### 8.6 Patterns to implement

- Path `:id` parsing: non-numeric, zero, or negative → 400.
- `created_by_id` / `updated_by_id` are set from the authenticated user, never from the body.
- Duplicate checks: rely on the DB unique index and translate the driver error into 409 (also pre-check in service for a clean message; handle the race by also mapping the DB error).
- Normalize: trim strings; lowercase emails.
- Limit request body to 1 MB (`fiber.Config{BodyLimit: 1 << 20}`).
- Reject non-JSON content type on write endpoints with 400 (or let the body parser fail → 400).

---

## 9. Authentication design

**JWT claims** (HS256):
```go
type Claims struct {
    UserID uint   `json:"uid"`
    Role   string `json:"role"`
    jwt.RegisteredClaims // sub=userID string, jti=uuid, iat, exp, iss="icecream-inventory"
}
```
- Generate `jti` with `github.com/google/uuid` (add to approved deps) or `crypto/rand` hex — prefer `crypto/rand` hex to stay dependency-free.
- On parse: verify signing method is HMAC (reject `none`/RS*), check `iss`, check `exp`.
- **Authentication middleware**: read `Authorization: Bearer <token>`; parse; check `jti` is not in `revoked_tokens`; load user from DB by `uid` (so deleted/changed users are handled — role comes from the **DB**, not from the token claim, so role changes take effect immediately); store `userID`, `role`, `jti`, `exp` in `c.Locals`.
- **Logout**: insert `jti` + `exp` into `revoked_tokens`; idempotent. Add a small background cleanup (delete rows where `expires_at < now`) on a ticker started in `main.go`, or run cleanup opportunistically on logout.
- Passwords: bcrypt; enforce min length 8 only on seed/creation helpers (no user-creation endpoint exists).

---

## 10. Security & reliability checklist

- `helmet`, `recover`, `requestid`, structured `logger` middleware enabled.
- CORS restricted to `CORS_ALLOW_ORIGINS`.
- Login rate limiter (`limiter` middleware on the login route only).
- No password/token/secret in logs. Redact `Authorization` header.
- Central Fiber `ErrorHandler`: maps `*httpx.AppError` and `*fiber.Error` to the envelope; everything else → 500 with generic message (log the real error with request id).
- GORM logger: `Warn` level in production, `Info` only when `APP_ENV=development`.
- Graceful shutdown: listen for SIGINT/SIGTERM, `app.ShutdownWithTimeout(10s)`, close DB.
- DB pool settings: `SetMaxOpenConns(25)`, `SetMaxIdleConns(5)`, `SetConnMaxLifetime(30m)` (SQLite in-memory: `SetMaxOpenConns(1)`).
- Use parameterized queries only (GORM does this; never string-concatenate user input into SQL). For `search`, use `LOWER(name) LIKE ?` with `%term%` and escape `%`/`_` in the term.

---

## 11. Seeding

`internal/database/seed.go` — **idempotent** (use `FirstOrCreate` on email / name; running twice must not duplicate or error):

- 3 users (credentials from env, defaults in section 6): Owner, Manager, Cashier.
- ~10 sample items: Chocolate, Vanilla, Strawberry, Mint Chip, Butterscotch (category `flavor`, unit `tub`), plus Sprinkles (`topping`), Waffle Cone (`cone`), Sugar Cone, Paper Cup (`cup`), Cone Sleeve (`packaging`). Realistic cents prices and one item with `quantity <= reorder_level` to demonstrate `low_stock`.
- ~5 sample customers.
- Seeding of sample items/customers only when `APP_ENV != production`; seeding of the three role users always (they are required by the deliverables).
- Entry point: `go run ./cmd/seed` (also `make seed`).

---

## 12. Phased execution plan (with acceptance criteria)

Track every task in `PROGRESS.md` (format in section 15). Time budget is a guide.

### Phase 0 — Bootstrap (≈ 45 min)
- `git init`, `.gitignore` (binaries, `.env`, `*.out`, `coverage.*`, IDE folders), `go mod init github.com/MrVirul/icecream-inventory` (change if the owner's module path differs).
- Create the directory skeleton, `.env.example`, `Makefile`, `docker-compose.yml` (postgres:16-alpine with healthcheck + named volume + the API service optional), `Dockerfile` (multi-stage, non-root user, static binary).
- Create `PROGRESS.md` and `docs/DECISIONS.md`.
- **Accept:** `go build ./...` passes; `docker compose up -d db` starts Postgres.

### Phase 1 — Core infrastructure (≈ 1.5 h)
- `config`, `database` (Postgres + SQLite openers), `models`, `migrate`, `httpx` (envelope, `AppError`, pagination parser, validator wrapper that maps `validator.ValidationErrors` to `details`), central error handler.
- `cmd/migrate/main.go`.
- **Accept:** `make migrate` creates all 4 tables in Postgres; unit tests for config validation and pagination parsing pass.

### Phase 2 — Authentication (≈ 2 h)
- `auth/password.go`, `auth/jwt.go`, `auth/permissions.go`.
- Repositories: users, revoked tokens. Service: login, logout.
- Middleware: `Authenticate`, `RequirePermission(p)`.
- Handlers + routes: `/api/auth/login`, `/api/auth/logout`, `/healthz`.
- **Accept:** login returns a token for each seeded role; wrong password → 401; logout then reuse token → 401; tampered/expired token → 401.

### Phase 3 — Inventory (≈ 2.5 h)
- DTOs, repository (list with search/filter/sort/pagination), service, handlers, routes with permission middleware exactly per section 8.4.
- **Accept:** every endpoint works with the Owner token; Manager/Cashier restrictions return 403 per the matrix; validation and conflict paths return 422/409.

### Phase 4 — Customers (≈ 2 h)
- Same pattern as inventory (copy the structure, don't invent a new one).
- **Accept:** same criteria as Phase 3 for the customer matrix.

### Phase 5 — Seeds (≈ 45 min)
- Section 11. `cmd/seed`, `make seed`.
- **Accept:** `make reset-db && make migrate && make seed` yields working logins; running `make seed` twice is a no-op.

### Phase 6 — Tests (≈ 4 h) — see section 14
- **Accept:** `go test -race ./...` green; coverage ≥ 80 % on `service`, `middleware`, `auth`; the RBAC matrix test covers all 8 operations × 3 roles + unauthenticated.

### Phase 7 — Documentation & diagrams (≈ 2 h)
- `docs/ER_DIAGRAM.md` (section 7.3), `docs/API.md` (section 16), `docs/openapi.yaml`, `docs/TESTING.md`, `README.md`.
- **Accept:** a new developer can clone, follow the README, and hit every endpoint using the docs alone. Every example `curl` in the docs was actually executed against the running server by the agent and the real output matches.

### Phase 8 — Hardening & final audit (≈ 2 h)
- Run the Final Audit (section 18). Fix everything found. Tag the release commit.

**Suggested schedule:** Day 1 = Phases 0–5 (working API). Day 2 = Phases 6–8 (tests, docs, audit) with a 2-hour buffer.

---

## 13. Verification Gate (run after every phase)

```bash
gofmt -l .            # must print nothing
go vet ./...
go build ./...
go test -race -count=1 ./...
```

From Phase 2 onward, also smoke-test against a **running** server:

```bash
make run &            # or: docker compose up -d
TOKEN=$(curl -s -X POST localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"owner@icecream.local","password":"Owner@12345"}' | jq -r .data.access_token)
curl -s localhost:8080/api/items -H "Authorization: Bearer $TOKEN" | jq
```

If the gate fails: fix before continuing. Do not "park" red tests.

---

## 14. Testing strategy

**Principles:** tests are fast (<15 s total), hermetic (in-memory SQLite, no network, no Docker), and deterministic. Use table-driven tests and `testify/require`.

### 14.1 Test utilities (`internal/testutil`)
- `NewTestApp(t) *TestApp` — opens in-memory SQLite (`file:<unique>?mode=memory&cache=shared`, `SetMaxOpenConns(1)`), runs migrations, seeds the 3 users (bcrypt cost 4), builds the real Fiber app via `router.New(...)`, returns helpers.
- `app.Token(t, role)` — logs in through the real endpoint and returns a bearer token.
- `app.Do(t, method, path, token, body) (*http.Response, []byte)` — wraps `app.Test(req)` with JSON encoding; decode helpers for the envelope.

### 14.2 Required test suites

1. **`auth` unit tests** — password hash/compare; JWT issue/parse; expired token; wrong secret; `alg=none` rejected; wrong issuer rejected; `Can()` for every role × permission (table of 24 cases).
2. **`config` tests** — missing/short `JWT_SECRET`, defaults, invalid port.
3. **`httpx` tests** — pagination defaults/limits/invalid values; sort whitelist; LIKE escaping; validator error mapping.
4. **Service unit tests** (with in-memory fakes or SQLite) — item/customer create, duplicate → conflict, not found, partial update only touches provided fields, `low_stock` computation, email lowercase normalization.
5. **Auth integration tests** — login success (each role), wrong password, unknown email (identical response), malformed body, logout revokes token, logout idempotency, missing/malformed `Authorization` header, tampered token, deleted user's token rejected.
6. **RBAC matrix test (the most important one)** — a single table-driven test:

```go
// 8 operations × 3 roles + 8 × "no token"
{"items list",     GET,    "/api/items",        Owner:200, Manager:200, Cashier:200, Anon:401},
{"items create",   POST,   "/api/items",        Owner:201, Manager:201, Cashier:403, Anon:401},
{"items update",   PATCH,  "/api/items/{id}",   Owner:200, Manager:200, Cashier:200, Anon:401},
{"items delete",   DELETE, "/api/items/{id}",   Owner:204, Manager:403, Cashier:403, Anon:401},
{"customers list", GET,    "/api/customers",    Owner:200, Manager:200, Cashier:200, Anon:401},
{"customers create",POST,  "/api/customers",    Owner:201, Manager:201, Cashier:201, Anon:401},
{"customers update",PATCH, "/api/customers/{id}",Owner:200,Manager:200, Cashier:200, Anon:401},
{"customers delete",DELETE,"/api/customers/{id}",Owner:204,Manager:403, Cashier:403, Anon:401},
```
   Each case creates its own fixture row first so deletes don't interfere. Also assert PUT behaves like PATCH for permissions, and that a forbidden delete leaves the row intact (re-fetch → 200).
7. **Inventory API tests** — list pagination/search/category/low_stock/sort; get by id (200/404/400 for `abc`); create validations (each rule → 422 with correct `details.field`); duplicate name (case-insensitive) → 409; PUT requires all fields; PATCH with empty body → 400/422; delete then get → 404; `created_by_id`/`updated_by_id` set from token and **ignored if sent in body**.
8. **Customer API tests** — same as 7 adapted (phone format, optional unique email, search by phone/email/name).
9. **Check-order tests** — Cashier `DELETE /api/items/999999` → 403; Owner same → 404; invalid body from forbidden role → 403 (not 422).
10. **Health/middleware tests** — `/healthz` 200; panic recovered → 500 envelope; body > 1 MB → 413 (or the configured error mapped to the envelope).
11. **Seeder test** — run seed twice, counts unchanged; roles correct; passwords verify.

### 14.3 Commands
```bash
go test -race -count=1 ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```
Document the results and the coverage matrix in `docs/TESTING.md`.

---

## 15. `PROGRESS.md` (create in Phase 0, update constantly)

This file is the **handoff contract** between agents/tools. Keep it exactly in this shape:

```markdown
# PROGRESS

## Current status
Phase: <N> — <name>   |   Last gate: PASS/FAIL (<date-time>)
Next action: <one concrete sentence an agent can start on immediately>

## Phases
- [x] Phase 0 — Bootstrap
- [ ] Phase 1 — Core infrastructure
  - [x] config loader + tests
  - [ ] httpx envelope + error handler
...

## Decisions (links to docs/DECISIONS.md)

## Blockers
- <date> <what failed> <what was tried> <workaround chosen>

## How to run
<exact commands that work right now>
```

Rules: update it **before every commit**; tick only what passed the gate; the "Next action" line must always be accurate.

---

## 16. Documentation deliverables (detail)

### `docs/API.md`
For **each** endpoint: method + path, purpose, required role(s) (cross-reference the matrix), headers, path/query params, request body schema with example, success response with example, all error responses (status + `code` + example), and a runnable `curl` example. Include sections: Overview, Authentication flow, Roles & permissions matrix, Conventions (envelope, pagination, sorting, errors), Endpoints, Seed credentials (dev only), Health.

### `docs/openapi.yaml`
OpenAPI 3.0.3 with `bearerAuth` security scheme, reusable `components/schemas` (Item, ItemCreate, ItemPatch, Customer…, Error, Pagination meta), reusable error responses, and per-operation `x-required-roles`. Validate it by loading with any YAML parser at minimum.

### `README.md`
Overview, features, stack, prerequisites (Go, Docker), quick start (copy `.env.example`, `docker compose up -d db`, `make migrate`, `make seed`, `make run`), test commands, project structure, seeded accounts table, link to docs, design decisions summary.

### `Makefile` targets
`run`, `build`, `test`, `test-cover`, `lint` (`go vet` + `gofmt -l`), `migrate`, `seed`, `db-up`, `db-down`, `reset-db`, `docker-build`, `tidy`.

---

## 17. Git workflow

- Branch: work on `main` (solo, short deadline) — but commit often.
- Conventional commits: `feat(items): add list endpoint with search`, `test(rbac): add permission matrix`, `docs: add ER diagram`, `chore: ...`, `fix: ...`.
- One commit per completed task group; never commit failing builds. Update `PROGRESS.md` in the same commit.
- Final: tag `v1.0.0`.

---

## 18. Final Audit (Phase 8 checklist — every box must be verified, not assumed)

**Functional**
- [ ] All 12 endpoints in section 8.4 + `/healthz` respond as documented.
- [ ] RBAC matrix verified by the automated test **and** manually via `curl` for Cashier delete and Cashier create-item (both 403).
- [ ] Seed users can log in; seed is idempotent.
- [ ] Logout invalidates the token.

**Quality**
- [ ] `gofmt -l .` empty; `go vet ./...` clean; `go test -race ./...` green.
- [ ] No `TODO`/`FIXME`/debug prints (`grep -rn "TODO\|FIXME\|fmt.Println" --include=*.go .`).
- [ ] No secrets in the repo (`git grep -n "JWT_SECRET="` only matches `.env.example`).
- [ ] Handlers contain no GORM calls; services contain no Fiber imports.
- [ ] Error responses never leak internals.

**Deliverables**
- [ ] ER diagram matches the actual models (compare column by column).
- [ ] `docs/API.md` and `docs/openapi.yaml` match actual behavior (spot-check 5 random endpoints with curl).
- [ ] `README` quick start works from a clean clone (`docker compose down -v` then follow it verbatim).
- [ ] `docker compose up --build` brings up DB + API and serves `/healthz`.

---

## 19. Reference snippets (follow these shapes)

**Central error type**
```go
// internal/httpx/errors.go
type AppError struct {
    Status  int
    Code    string
    Message string
    Details []FieldError // optional
}
func (e *AppError) Error() string { return e.Message }

func BadRequest(msg string) *AppError { return &AppError{400, "BAD_REQUEST", msg, nil} }
func Unauthorized(msg string) *AppError { return &AppError{401, "UNAUTHORIZED", msg, nil} }
func Forbidden() *AppError { return &AppError{403, "FORBIDDEN", "you do not have permission to perform this action", nil} }
func NotFound(what string) *AppError { return &AppError{404, "NOT_FOUND", what + " not found", nil} }
func Conflict(msg string) *AppError { return &AppError{409, "CONFLICT", msg, nil} }
```

**RBAC middleware**
```go
func RequirePermission(p auth.Permission) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals(CtxRole).(string)
        if !ok || !auth.Can(auth.Role(role), p) {
            return httpx.Forbidden()
        }
        return c.Next()
    }
}
```

**Route wiring**
```go
api := app.Group("/api")
api.Post("/auth/login", loginLimiter, authH.Login)

protected := api.Group("", mw.Authenticate(authSvc))
protected.Post("/auth/logout", authH.Logout)

items := protected.Group("/items")
items.Get("/",    mw.RequirePermission(auth.ItemsRead),   itemH.List)
items.Get("/:id", mw.RequirePermission(auth.ItemsRead),   itemH.Get)
items.Post("/",   mw.RequirePermission(auth.ItemsCreate), itemH.Create)
items.Put("/:id", mw.RequirePermission(auth.ItemsUpdate), itemH.Replace)
items.Patch("/:id", mw.RequirePermission(auth.ItemsUpdate), itemH.Patch)
items.Delete("/:id", mw.RequirePermission(auth.ItemsDelete), itemH.Delete)
// customers: identical shape with Customers* permissions
```

**Router constructor must accept dependencies** (`router.New(cfg, db) *fiber.App`) so tests build the same app production uses.

**Item response shape**
```json
{
  "id": 1,
  "name": "Chocolate",
  "category": "flavor",
  "description": "Rich dark chocolate",
  "quantity": 12,
  "unit": "tub",
  "price_cents": 1850,
  "reorder_level": 5,
  "low_stock": false,
  "created_by_id": 1,
  "updated_by_id": 1,
  "created_at": "2026-10-27T10:00:00Z",
  "updated_at": "2026-10-27T10:00:00Z"
}
```

---

## 20. Prompts for the human to paste

**Kickoff (first session, opencode):**
> Read `agent.md` completely. Then create `PROGRESS.md` as specified in section 15 and start Phase 0. Work autonomously through the phases in order. After each phase, run the Verification Gate (section 13), update `PROGRESS.md`, and commit. Do not ask me questions — make the simplest reasonable decision and log it in `docs/DECISIONS.md`. Do not stop until Phase 8 is complete or you hit a usage limit; before stopping for any reason, make sure `PROGRESS.md` "Next action" is accurate and everything is committed.

**Resume (new session / switching to Antigravity):**
> Read `agent.md` completely, then read `PROGRESS.md` and `git log --oneline -20`. Run the Verification Gate. If it is red, fix that first. Then continue from the "Next action" line in `PROGRESS.md` and keep following the phases in `agent.md`. Same rules: no questions, log decisions, update `PROGRESS.md` before every commit.

**Focused fix prompts (use when something breaks):**
> The gate is failing. Run it, read the first failure only, fix the root cause (not the symptom), re-run. Repeat until green. Do not modify tests to make them pass unless the test contradicts `agent.md`.

> Audit mode: execute section 18 of `agent.md` item by item. For each box, show the command you ran and its output. Fix every failure and re-verify.

---

## 21. Definition of Done

The project is done when **all** are true:

1. Section 18 audit is fully checked.
2. A clean clone + README quick start yields a working API with seeded Owner/Manager/Cashier logins.
3. `go test -race ./...` is green, including the RBAC matrix.
4. ER diagram, API docs, and OpenAPI spec exist and match the running behavior.
5. `PROGRESS.md` shows every phase ticked and the repo is tagged `v1.0.0`.

*End of agent.md.*
