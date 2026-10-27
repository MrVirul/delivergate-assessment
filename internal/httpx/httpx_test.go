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
