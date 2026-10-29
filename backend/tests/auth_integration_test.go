package tests

import (
	"net/http"
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
	ta := testutil.NewTestApp(t)
	// wrong password
	resp, body := ta.Do(t, http.MethodPost, "/api/auth/login", "", map[string]string{"email": "owner@icecream.local", "password": "wrong"})
	require.Equal(t, 401, resp.StatusCode)
	var res map[string]interface{}
	require.NoError(t, decodeJSON(body, &res))
	errObj := res["error"].(map[string]interface{})
	require.Equal(t, "invalid email or password", errObj["message"])

	// unknown email same message
	resp, body = ta.Do(t, http.MethodPost, "/api/auth/login", "", map[string]string{"email": "nope@x.com", "password": "wrongpass123"})
	require.Equal(t, 401, resp.StatusCode)
	require.NoError(t, decodeJSON(body, &res))
	errObj = res["error"].(map[string]interface{})
	require.Equal(t, "invalid email or password", errObj["message"])

	// logout
	token := ta.Token(t, "Owner")
	resp, _ = ta.Do(t, http.MethodPost, "/api/auth/logout", token, nil)
	require.Equal(t, 204, resp.StatusCode)
	// reuse
	resp, _ = ta.Do(t, http.MethodGet, "/api/items", token, nil)
	require.Equal(t, 401, resp.StatusCode)
}
