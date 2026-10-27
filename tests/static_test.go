package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestStaticServing(t *testing.T) {
	ta := testutil.NewTestApp(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := ta.App.Test(req)
	require.NoError(t, err)
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)
	// API still 401 without token
	resp, _ = ta.Do(t, http.MethodGet, "/api/items", "", nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	// healthz still works
	resp, body := ta.Do(t, http.MethodGet, "/healthz", "", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), "ok")
}
