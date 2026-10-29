package auth_test

import (
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashCompare(t *testing.T) {
	h, err := auth.HashPassword("testpass123", 4)
	require.NoError(t, err)
	require.NoError(t, auth.ComparePassword(h, "testpass123"))
	require.Error(t, auth.ComparePassword(h, "wrong"))
}

func TestPermissions(t *testing.T) {
	require.True(t, auth.Can(auth.RoleOwner, auth.ItemsDelete))
	require.False(t, auth.Can(auth.RoleManager, auth.ItemsDelete))
	require.False(t, auth.Can(auth.RoleCashier, auth.ItemsDelete))
	require.True(t, auth.Can(auth.RoleManager, auth.ItemsCreate))
	require.False(t, auth.Can(auth.RoleCashier, auth.ItemsCreate))
	require.True(t, auth.Can(auth.RoleCashier, auth.ItemsRead))
	require.True(t, auth.Can(auth.RoleCashier, auth.CustomersCreate))
	require.False(t, auth.Can(auth.Role(""), auth.ItemsRead))
}
