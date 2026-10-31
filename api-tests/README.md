# API Tests (Bruno)

A [Bruno](https://www.usebruno.com/) collection for the Ice Cream Inventory API.

## Open the collection

In Bruno, choose **Open Collection** and select this `api-tests/` directory
(the folder containing `bruno.json`).

Then select the **Local** environment from the environment dropdown.

## Run the collection

Run the whole collection top to bottom. The requests are ordered so that:

1. `Auth > Login` obtains a JWT and stores it in the `token` runtime variable.
2. `Health`, `Items`, and `Customers` use that token via collection-level
   Bearer auth.
3. `Create Item` / `Create Customer` store the new IDs in `itemId` /
   `customerId` for the follow-up requests.
4. `Cleanup > Logout` revokes the token last.

You can also run requests individually, but run **Login** first so the token is
available.

## Run from the CLI

```bash
# install once
npm install -g @usebruno/cli

# run from the collection root; -r recurses into folders
cd api-tests
bru run . -r --env Local
```

Override the server URL without editing files:

```bash
bru run . -r --env Local --env-var baseUrl=http://localhost:9090
```

Without `-r`, the CLI only runs requests at the collection root, which is why
the command above includes it.

## Configuration

| Variable          | Default                 | Purpose                     |
| ----------------- | ----------------------- | --------------------------- |
| `baseUrl`         | `http://localhost:8080` | Backend base URL            |
| `ownerEmail`      | `owner@mail.com`        | Seed Owner login            |
| `ownerPassword`   | `Owner@123`             | Seed Owner password         |
| `managerEmail`    | `manager@mail.com`      | Seed Manager login          |
| `managerPassword` | `Manager@123`           | Seed Manager password       |
| `cashierEmail`    | `cashier@mail.com`      | Seed Cashier login          |
| `cashierPassword` | `Cashier@123`           | Seed Cashier password       |

The credentials match `backend/.env`. If you override the `SEED_*` variables,
update `environments/Local.bru`.

Other environments (staging, production) can be added under `environments/`;
copy `Local.bru` and change `baseUrl` and credentials.

## Saved response examples

Every request ships with one saved example under its **Examples** tab in Bruno.
The examples use realistic shop data and keep the full response body, headers,
and status code, so you can review the API contract without running the server.

They were captured against the `Local` environment:

- Items use the seeded Sri Lankan ice cream catalogue (`Ceylon Vanilla`,
  `Kithul Treacle`, `Wood Apple`, `King Coconut`, `Mango`, `Curd and Honey`,
  plus toppings, cones, cups, and packaging).
- Customers use the seeded Colombo/Kandy/Galle/Negombo/Jaffna customer book.
- The write examples follow one realistic run: create `Ceylon Vanilla Reserve`,
  restock/reprice it, mark it low stock, then delete it; create customer
  `Ishara Weerasinghe`, update her address after a move, then delete her.

If you change a request body, update its saved example too so the documented
request and response stay in sync.

## Test data

- **Item:** `Ceylon Vanilla Reserve` (flavor, 18 tubs at LKR 1,850.00) is created,
  replaced with a restocked/repriced version, patched to 7 tubs, and deleted.
- **Customer:** `Ishara Weerasinghe` (`ishara.weerasinghe@gmail.com`,
  `+94 77 987 6543`, Flower Road, Colombo 07) is created, moved to Havelock
  Road, patched, and deleted.

The full run cleans up both records and logs out, so it can be repeated against
the same database. A run that fails midway can leave the item or customer
behind; delete it in the app (or Postman/Bruno) before running again.

## Coverage

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /healthz`
- `GET /api/items`, `POST /api/items`
- `GET`, `PUT`, `PATCH`, `DELETE /api/items/:id`
- `GET /api/customers`, `POST /api/customers`
- `GET`, `PUT`, `PATCH`, `DELETE /api/customers/:id`

Each request includes status/body assertions, and the create/get/replace/patch
requests assert the realistic values written by the previous step.
