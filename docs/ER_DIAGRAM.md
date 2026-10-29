# ER Diagram

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
        datetime created_at
        datetime updated_at
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
        datetime created_at
        datetime updated_at
    }
    REVOKED_TOKENS {
        uint id PK
        string jti UK
        datetime expires_at
        datetime created_at
    }
```

## Entity Descriptions

### users
Stores authenticated users with roles (Owner, Manager, Cashier). Passwords are stored as bcrypt hashes.

### inventory_items
Inventory items with pricing in cents, stock levels, reorder levels, and computed low_stock status.

### customers
Customer records with optional unique email and required phone.

### revoked_tokens
Stores revoked JWT JTIs for logout functionality.
