# API Documentation

## Overview
Base path: `/api`. All bodies are JSON. Content-Type `application/json`. Responses use envelope format.

### Success Envelope
- Single: `{"data": { ... }}`
- List: `{"data": [ ... ], "meta": {"page": 1, "limit": 20, "total": 57, "total_pages": 3}}`

### Error Envelope
`{"error": {"code": "VALIDATION_ERROR", "message": "...", "details": [{"field": "...", "message": "..."}]}}`

## Authentication
### POST /api/auth/login (public)
Body: `{"email": "...", "password": "..."}`
Response: `{"data": {"access_token": "...", "token_type": "Bearer", "expires_in": 3600, "user": {"id": 1, "name": "...", "email": "...", "role": "Owner"}}}`

### POST /api/auth/logout (authenticated)
Headers: `Authorization: Bearer <token>`
Response: `204 No Content` (revokes token)

## Roles & Permissions
See matrix in agent.md section 4.

## Endpoints

### Health
#### GET /healthz (public)
Response: `{"status": "ok"}` (503 if DB down)

### Inventory
#### GET /api/items (items:read)
Query: `search`, `category`, `low_stock=true`, `page`, `limit`, `sort` (name|quantity|price_cents|created_at, prefix `-` for desc)
#### GET /api/items/:id (items:read)
#### POST /api/items (items:create)
#### PUT /api/items/:id (items:update)
#### PATCH /api/items/:id (items:update)
#### DELETE /api/items/:id (items:delete)

### Customers
#### GET /api/customers (customers:read)
Query: `search`, `page`, `limit`, `sort` (name|created_at)
#### GET /api/customers/:id (customers:read)
#### POST /api/customers (customers:create)
#### PUT /api/customers/:id (customers:update)
#### PATCH /api/customers/:id (customers:update)
#### DELETE /api/customers/:id (customers:delete)

## Seed Credentials (dev)
The defaults below come from [`backend/.env.example`](../backend/.env.example). If you override the `SEED_*` variables in your `.env`, use those values instead.
- Owner: owner@mail.com / Owner@123
- Manager: manager@mail.com / Manager@123
- Cashier: cashier@mail.com / Cashier@123
