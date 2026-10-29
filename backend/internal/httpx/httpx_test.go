package httpx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginationDefaults(t *testing.T) {
	// just test helper functions
	require.Equal(t, 2, TotalPages(25, 20))
	require.Equal(t, 0, TotalPages(0, 20))
}

func TestValidateStructUsesJSONFieldNames(t *testing.T) {
	type payload struct {
		Name     string `json:"name" validate:"required,min=2"`
		Price    *int64 `json:"price_cents" validate:"required,gte=0"`
		Optional string `json:"optional,omitempty" validate:"max=5"`
	}

	err := ValidateStruct(payload{Name: ""})
	require.Error(t, err)

	appErr, ok := err.(*AppError)
	require.True(t, ok)
	require.Equal(t, 422, appErr.Status)
	require.NotEmpty(t, appErr.Details)

	fields := make(map[string]string, len(appErr.Details))
	for _, d := range appErr.Details {
		fields[d.Field] = d.Message
	}
	require.Contains(t, fields, "name")
	require.Contains(t, fields, "price_cents")
	require.NotContains(t, fields, "PriceCents")
}

