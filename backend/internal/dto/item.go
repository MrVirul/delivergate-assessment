package dto

type ItemCreateRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=120"`
	Category     string `json:"category" validate:"required,oneof=flavor topping cone cup packaging other"`
	Description  string `json:"description" validate:"max=1000"`
	Quantity     *int   `json:"quantity" validate:"required,gte=0"`
	Unit         string `json:"unit" validate:"required,oneof=tub scoop kg litre pcs box"`
	PriceCents   *int64 `json:"price_cents" validate:"required,gte=0"`
	ReorderLevel *int   `json:"reorder_level" validate:"omitempty,gte=0"`
}

type ItemUpdateRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=120"`
	Category     string `json:"category" validate:"required,oneof=flavor topping cone cup packaging other"`
	Description  string `json:"description" validate:"max=1000"`
	Quantity     *int   `json:"quantity" validate:"required,gte=0"`
	Unit         string `json:"unit" validate:"required,oneof=tub scoop kg litre pcs box"`
	PriceCents   *int64 `json:"price_cents" validate:"required,gte=0"`
	ReorderLevel *int   `json:"reorder_level" validate:"omitempty,gte=0"`
}

type ItemPatchRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=2,max=120"`
	Category     *string `json:"category" validate:"omitempty,oneof=flavor topping cone cup packaging other"`
	Description  *string `json:"description" validate:"omitempty,max=1000"`
	Quantity     *int    `json:"quantity" validate:"omitempty,gte=0"`
	Unit         *string `json:"unit" validate:"omitempty,oneof=tub scoop kg litre pcs box"`
	PriceCents   *int64  `json:"price_cents" validate:"omitempty,gte=0"`
	ReorderLevel *int    `json:"reorder_level" validate:"omitempty,gte=0"`
}

type ItemResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	Quantity     int    `json:"quantity"`
	Unit         string `json:"unit"`
	PriceCents   int64  `json:"price_cents"`
	ReorderLevel int    `json:"reorder_level"`
	LowStock     bool   `json:"low_stock"`
	CreatedByID  *uint  `json:"created_by_id"`
	UpdatedByID  *uint  `json:"updated_by_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
