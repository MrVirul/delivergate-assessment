package dto

type CustomerCreateRequest struct {
	Name    string  `json:"name" validate:"required,min=2,max=120"`
	Email   *string `json:"email" validate:"omitempty,email"`
	Phone   string  `json:"phone" validate:"required,min=7,max=30"`
	Address string  `json:"address" validate:"max=255"`
	Notes   string  `json:"notes" validate:"max=1000"`
}

type CustomerUpdateRequest struct {
	Name    string  `json:"name" validate:"required,min=2,max=120"`
	Email   *string `json:"email" validate:"omitempty,email"`
	Phone   string  `json:"phone" validate:"required,min=7,max=30"`
	Address string  `json:"address" validate:"max=255"`
	Notes   string  `json:"notes" validate:"max=1000"`
}

type CustomerPatchRequest struct {
	Name    *string `json:"name" validate:"omitempty,min=2,max=120"`
	Email   *string `json:"email" validate:"omitempty,email"`
	Phone   *string `json:"phone" validate:"omitempty,min=7,max=30"`
	Address *string `json:"address" validate:"omitempty,max=255"`
	Notes   *string `json:"notes" validate:"omitempty,max=1000"`
}

type CustomerResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Email       *string `json:"email"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	Notes       string  `json:"notes"`
	CreatedByID *uint   `json:"created_by_id"`
	UpdatedByID *uint   `json:"updated_by_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
