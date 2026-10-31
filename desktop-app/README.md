# Desktop App (desktop-app)

WPF (.NET 8) desktop client for the Ice Cream Shop Inventory API at `http://localhost:8080`. It implements the login, inventory, customer, and role-permission screens from `doc.md` against the real backend contract.

## Prerequisites

- Windows 10/11 (WPF is Windows-only)
- [.NET 8 SDK](https://dotnet.microsoft.com/download/dotnet/8.0)
- A running backend with `make migrate && make seed` already completed

## Project layout

| Path                              | Purpose                                                    |
| --------------------------------- | ---------------------------------------------------------- |
| `Views/LoginWindow.xaml`          | Email/password login with validation and error banner      |
| `ViewModels/LoginViewModel.cs`    | Login flow and error handling                              |
| `MainWindow.xaml`                 | Dashboard with Inventory and Customer tabs                 |
| `ViewModels/MainViewModel.cs`     | Data loading, add/edit commands, RBAC checks, error banner |
| `Views/ItemEditorWindow.xaml`     | Add/Edit inventory item dialog                             |
| `Views/CustomerEditorWindow.xaml` | Add/Edit customer dialog                                   |
| `Services/ApiService.cs`          | HTTP client for the backend API                            |
| `Models/`                         | Items, customers, login and API envelope DTOs              |

## Running

Start the backend first (see [`../backend/README.md`](../backend/README.md)):

```bash
cd ../backend
cp .env.example .env
# Edit .env and set DATABASE_URL / JWT_SECRET
make migrate
make seed
make run
```

Then run the desktop client from this directory:

```bash
dotnet restore
dotnet run
```

Alternatively, from the repository root:

```bash
dotnet run --project desktop-app
```

The client defaults to `http://localhost:8080/api/` even when no configuration
file is present, so it works out of the box against the local backend.

## Configuration

The API base URL and page size are read from the optional `appsettings.json`,
then the optional `appsettings.Local.json`, then environment variables. If none
of these are present, the built-in defaults below are used:

```json
{
  "Api": {
    "BaseUrl": "http://localhost:8080/api/",
    "PageSize": 100
  }
}
```

To override them locally without touching tracked files, copy the example and edit it:

```bash
cp appsettings.Local.example.json appsettings.Local.json
```

`appsettings.Local.json` is gitignored. Environment variables take precedence over both JSON files:

```bash
DESKTOPAPP_Api__BaseUrl=http://localhost:9000/api/ dotnet run
DESKTOPAPP_Api__PageSize=50 dotnet run
```

A `.env` file is not read by .NET out of the box; use the JSON override or an
environment variable (or add the `DotNetEnv` package if you prefer `.env`).

## Seeded login accounts

These are the defaults created by `make seed` from [`backend/.env.example`](../backend/.env.example):

| Role    | Email              | Password      |
| ------- | ------------------ | ------------- |
| Owner   | `owner@mail.com`   | `Owner@123`   |
| Manager | `manager@mail.com` | `Manager@123` |
| Cashier | `cashier@mail.com` | `Cashier@123` |

If the backend's `.env` overrides the `SEED_*` variables, use those credentials instead.

## Backend endpoints used

| Method   | Endpoint                               | Purpose                               |
| -------- | -------------------------------------- | ------------------------------------- |
| `POST`   | `auth/login`                           | Email/password login                  |
| `POST`   | `auth/logout`                          | Revoke the current token              |
| `GET`    | `items?page=1&limit=100&sort=name`     | List inventory (all pages are loaded) |
| `POST`   | `items`                                | Create an item                        |
| `PUT`    | `items/{id}`                           | Update an item                        |
| `DELETE` | `items/{id}`                           | Delete an item                        |
| `GET`    | `customers?page=1&limit=100&sort=name` | List customers (all pages are loaded) |
| `POST`   | `customers`                            | Create a customer                     |
| `PUT`    | `customers/{id}`                       | Update a customer                     |
| `DELETE` | `customers/{id}`                       | Delete a customer                     |

All responses use the backend envelope format (`{ "data": ... }`, list `meta`, and `{ "error": ... }`). Item prices are stored as `price_cents` and displayed/edited as dollars.

## Permissions

The UI mirrors the backend RBAC matrix:

| Capability            | Owner | Manager | Cashier |
| --------------------- | :---: | :-----: | :-----: |
| View inventory        |   ✓   |    ✓    |    ✓    |
| Add inventory item    |   ✓   |    ✓    |         |
| Edit inventory item   |   ✓   |    ✓    |    ✓    |
| Delete inventory item |   ✓   |         |         |
| View customers        |   ✓   |    ✓    |    ✓    |
| Add/edit customers    |   ✓   |    ✓    |    ✓    |
| Delete customer       |   ✓   |         |         |

If the API returns `403 Forbidden`, the dashboard displays:

`Action Error: You lack sufficient permissions to perform this operation.`

Server validation and conflict messages (for example, duplicate item names) are shown in the same banner.
