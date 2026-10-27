package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/testutil"
	"github.com/stretchr/testify/require"
)

func decodeJSON(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func TestRBACMatrix(t *testing.T) {
	ta := testutil.NewTestApp(t)

	// Create a fixture item
	itemBody := map[string]interface{}{
		"name":         "TestItem",
		"category":     "flavor",
		"quantity":     10,
		"unit":         "tub",
		"price_cents":  1000,
		"reorder_level": 5,
	}
	resp, body := ta.Do(t, http.MethodPost, "/api/items", ta.Token(t, "Owner"), itemBody)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var itemRes map[string]interface{}
	require.NoError(t, decodeJSON(body, &itemRes))
	data := itemRes["data"].(map[string]interface{})
	itemID := int(data["id"].(float64))

	// Create a fixture customer
	custBody := map[string]interface{}{
		"name":  "TestCust",
		"phone": "1234567890",
	}
	resp, body = ta.Do(t, http.MethodPost, "/api/customers", ta.Token(t, "Owner"), custBody)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var custRes map[string]interface{}
	require.NoError(t, decodeJSON(body, &custRes))
	cdata := custRes["data"].(map[string]interface{})
	custID := int(cdata["id"].(float64))

	// Items list
	resp, _ = ta.Do(t, http.MethodGet, "/api/items", ta.Token(t, "Owner"), nil)
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/items", ta.Token(t, "Manager"), nil)
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/items", ta.Token(t, "Cashier"), nil)
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/items", "", nil)
	require.Equal(t, 401, resp.StatusCode)

	// Items create
	resp, _ = ta.Do(t, http.MethodPost, "/api/items", ta.Token(t, "Manager"), map[string]interface{}{"name": "TestItem2", "category": "flavor", "quantity": 1, "unit": "tub", "price_cents": 500})
	require.Equal(t, 201, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPost, "/api/items", ta.Token(t, "Cashier"), map[string]interface{}{"name": "TestItem3", "category": "flavor", "quantity": 1, "unit": "tub", "price_cents": 500})
	require.Equal(t, 403, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPost, "/api/items", "", map[string]interface{}{"name": "TestItem4", "category": "flavor", "quantity": 1, "unit": "tub", "price_cents": 500})
	require.Equal(t, 401, resp.StatusCode)

	// Items update
	resp, _ = ta.Do(t, http.MethodPatch, "/api/items/"+itoa(itemID), ta.Token(t, "Owner"), map[string]interface{}{"quantity": 11})
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPatch, "/api/items/"+itoa(itemID), ta.Token(t, "Manager"), map[string]interface{}{"quantity": 12})
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPatch, "/api/items/"+itoa(itemID), ta.Token(t, "Cashier"), map[string]interface{}{"quantity": 13})
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPatch, "/api/items/"+itoa(itemID), "", map[string]interface{}{"quantity": 14})
	require.Equal(t, 401, resp.StatusCode)

	// Items delete - check order
	resp, _ = ta.Do(t, http.MethodDelete, "/api/items/"+itoa(999999), ta.Token(t, "Cashier"), nil)
	require.Equal(t, 403, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodDelete, "/api/items/"+itoa(itemID), ta.Token(t, "Manager"), nil)
	require.Equal(t, 403, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodDelete, "/api/items/"+itoa(itemID), ta.Token(t, "Cashier"), nil)
	require.Equal(t, 403, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodDelete, "/api/items/"+itoa(itemID), ta.Token(t, "Owner"), nil)
	require.Equal(t, 204, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/items/"+itoa(itemID), ta.Token(t, "Owner"), nil)
	require.Equal(t, 404, resp.StatusCode)

	// Customers
	resp, _ = ta.Do(t, http.MethodGet, "/api/customers", ta.Token(t, "Owner"), nil)
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/customers", ta.Token(t, "Manager"), nil)
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/customers", ta.Token(t, "Cashier"), nil)
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodGet, "/api/customers", "", nil)
	require.Equal(t, 401, resp.StatusCode)

	resp, _ = ta.Do(t, http.MethodPost, "/api/customers", ta.Token(t, "Owner"), map[string]interface{}{"name": "C2", "phone": "1111111111"})
	require.Equal(t, 201, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPost, "/api/customers", ta.Token(t, "Manager"), map[string]interface{}{"name": "C3", "phone": "2222222222"})
	require.Equal(t, 201, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPost, "/api/customers", ta.Token(t, "Cashier"), map[string]interface{}{"name": "C4", "phone": "3333333333"})
	require.Equal(t, 201, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPost, "/api/customers", "", map[string]interface{}{"name": "C5", "phone": "4444444444"})
	require.Equal(t, 401, resp.StatusCode)

	resp, _ = ta.Do(t, http.MethodPatch, "/api/customers/"+itoa(custID), ta.Token(t, "Owner"), map[string]interface{}{"name": "TestCustX"})
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPatch, "/api/customers/"+itoa(custID), ta.Token(t, "Manager"), map[string]interface{}{"name": "TestCustY"})
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPatch, "/api/customers/"+itoa(custID), ta.Token(t, "Cashier"), map[string]interface{}{"name": "TestCustZ"})
	require.Equal(t, 200, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodPatch, "/api/customers/"+itoa(custID), "", map[string]interface{}{"name": "X"})
	require.Equal(t, 401, resp.StatusCode)

	resp, _ = ta.Do(t, http.MethodDelete, "/api/customers/"+itoa(custID), ta.Token(t, "Manager"), nil)
	require.Equal(t, 403, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodDelete, "/api/customers/"+itoa(custID), ta.Token(t, "Cashier"), nil)
	require.Equal(t, 403, resp.StatusCode)
	resp, _ = ta.Do(t, http.MethodDelete, "/api/customers/"+itoa(custID), ta.Token(t, "Owner"), nil)
	require.Equal(t, 204, resp.StatusCode)
}
