package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MrVirul/icecream-inventory/internal/config"
	"github.com/MrVirul/icecream-inventory/internal/database"
	"github.com/MrVirul/icecream-inventory/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestApp struct {
	App *fiber.App
	DB  *gorm.DB
	Cfg *config.Config
}

func NewTestApp(t *testing.T) *TestApp {
	t.Helper()
	cfg := &config.Config{
		AppEnv:               "test",
		AppPort:              "8080",
		DatabaseURL:          "",
		JWTSecret:            "test-secret-key-that-is-at-least-32-chars!",
		JWTTTLMinutes:        60,
		AutoMigrate:          true,
		BcryptCost:           4,
		CORSAllowOrigins:     "*",
		LoginRateLimitPerMin: 1000,
		SeedOwnerEmail:       "owner@mail.com",
		SeedOwnerPassword:    "Owner@123",
		SeedManagerEmail:     "manager@mail.com",
		SeedManagerPassword:  "Manager@123",
		SeedCashierEmail:     "cashier@mail.com",
		SeedCashierPassword:  "Cashier@123",
	}
	db, err := database.OpenSQLite(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	require.NoError(t, err)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.Migrate(db))
	require.NoError(t, database.Seed(db, database.SeedConfig{
		AppEnv:              cfg.AppEnv,
		BcryptCost:          cfg.BcryptCost,
		SeedOwnerEmail:      cfg.SeedOwnerEmail,
		SeedOwnerPassword:   cfg.SeedOwnerPassword,
		SeedManagerEmail:    cfg.SeedManagerEmail,
		SeedManagerPassword: cfg.SeedManagerPassword,
		SeedCashierEmail:    cfg.SeedCashierEmail,
		SeedCashierPassword: cfg.SeedCashierPassword,
	}))

	app := router.New(cfg, db)
	return &TestApp{App: app, DB: db, Cfg: cfg}
}

func (ta *TestApp) Token(t *testing.T, role string) string {
	t.Helper()
	var email, pass string
	switch role {
	case "Owner":
		email = ta.Cfg.SeedOwnerEmail
		pass = ta.Cfg.SeedOwnerPassword
	case "Manager":
		email = ta.Cfg.SeedManagerEmail
		pass = ta.Cfg.SeedManagerPassword
	case "Cashier":
		email = ta.Cfg.SeedCashierEmail
		pass = ta.Cfg.SeedCashierPassword
	default:
		email = ta.Cfg.SeedOwnerEmail
		pass = ta.Cfg.SeedOwnerPassword
	}
	body := map[string]string{"email": email, "password": pass}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := ta.App.Test(req, -1)
	require.NoError(t, err)
	var res map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
	data := res["data"].(map[string]interface{})
	token, _ := data["access_token"].(string)
	return token
}

func (ta *TestApp) Do(t *testing.T, method, path, token string, body interface{}) (*http.Response, []byte) {
	t.Helper()
	var reqBody *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}
	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := ta.App.Test(req, -1)
	require.NoError(t, err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	_ = time.Now()
	return resp, buf.Bytes()
}
