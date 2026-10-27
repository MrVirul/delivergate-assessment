package auth

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

func set(ps ...Permission) map[Permission]struct{} {
	m := make(map[Permission]struct{}, len(ps))
	for _, p := range ps {
		m[p] = struct{}{}
	}
	return m
}

var rolePermissions = map[Role]map[Permission]struct{}{
	RoleOwner: set(ItemsRead, ItemsCreate, ItemsUpdate, ItemsDelete,
		CustomersRead, CustomersCreate, CustomersUpdate, CustomersDelete),
	RoleManager: set(ItemsRead, ItemsCreate, ItemsUpdate,
		CustomersRead, CustomersCreate, CustomersUpdate),
	RoleCashier: set(ItemsRead, ItemsUpdate,
		CustomersRead, CustomersCreate, CustomersUpdate),
}

func Can(role Role, p Permission) bool {
	if role == "" {
		return false
	}
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	_, ok = perms[p]
	return ok
}
