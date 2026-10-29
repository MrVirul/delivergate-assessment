# Desktop App (desktop-app)

WPF (.NET 8) desktop client for the Ice Cream Shop Inventory API at `http://localhost:8080`. It implements the login, inventory, customer, and role-permission screens from `doc.md` against the real backend contract.

## Project layout

| Path | Purpose |
| --- | --- |
| `Views/LoginWindow.xaml` | Email/password login with validation and error banner |
| `ViewModels/LoginViewModel.cs` | Login flow and error handling |
| `MainWindow.xaml` | Dashboard with Inventory and Customer tabs |
| `ViewModels/MainViewModel.cs` | Data loading, add/edit commands, RBAC checks, error banner |
| `Views/ItemEditorWindow.xaml` | Add/Edit inventory item dialog |
| `Views/CustomerEditorWindow.xaml` | Add/Edit customer dialog |
| `Services/ApiService.cs` | HTTP client for the backend API |
| `Models/` | Items, customers, login and API envelope DTOs |

## Running

Start the backend first:

```bash
cd ../backend
make migrate
make seed
make run
```

Then run the desktop client:

```bash
dotnet restore
dotnet run
```

The API base URL is `http://localhost:8080/api/` in `Services/ApiService.cs`.

## Seeded login accounts

These come from the backend `.env` file:

| Role | Email | Password |
| --- | --- | --- |
| Owner | `owner@mail.com` | `Owner@123` |
| Manager | `manager@mail.com` | `Manager@123` |
| Cashier | `cashier@mail.com` | `Cashier@123` |

## Backend endpoints used

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `auth/login` | Email/password login |
| `POST` | `auth/logout` | Revoke the current token |
| `GET` | `items?page=1&limit=100&sort=name` | List inventory (all pages are loaded) |
| `POST` | `items` | Create an item |
| `PUT` | `items/{id}` | Update an item |
| `GET` | `customers?page=1&limit=100&sort=name` | List customers (all pages are loaded) |
| `POST` | `customers` | Create a customer |
| `PUT` | `customers/{id}` | Update a customer |

All responses use the backend envelope format (`{ "data": ... }`, list `meta`, and `{ "error": ... }`). Item prices are stored as `price_cents` and displayed/edited as dollars.

## Permissions

The UI mirrors the backend RBAC matrix:

| Capability | Owner | Manager | Cashier |
| --- | :---: | :---: | :---: |
| View inventory | ✓ | ✓ | ✓ |
| Add inventory item | ✓ | ✓ | |
| Edit inventory item | ✓ | ✓ | ✓ |
| View customers | ✓ | ✓ | ✓ |
| Add/edit customers | ✓ | ✓ | ✓ |

If the API returns `403 Forbidden`, the dashboard displays:

`Action Error: You lack sufficient permissions to perform this operation.`

Server validation and conflict messages (for example, duplicate item names) are shown in the same banner.
