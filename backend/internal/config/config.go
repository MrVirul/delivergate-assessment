package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	AppPort              string
	DatabaseURL          string
	JWTSecret            string
	JWTTTLMinutes        int
	AutoMigrate          bool
	BcryptCost           int
	CORSAllowOrigins     string
	LoginRateLimitPerMin int

	SeedOwnerEmail      string
	SeedOwnerPassword   string
	SeedManagerEmail    string
	SeedManagerPassword string
	SeedCashierEmail    string
	SeedCashierPassword string
}

func Load() (*Config, error) {
	// Load .env if present (optional)
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		AppPort:              getEnv("APP_PORT", "8080"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		JWTTTLMinutes:        getEnvInt("JWT_TTL_MINUTES", 60),
		AutoMigrate:          getEnvBool("AUTO_MIGRATE", true),
		BcryptCost:           getEnvInt("BCRYPT_COST", 12),
		CORSAllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000"),
		LoginRateLimitPerMin: getEnvInt("LOGIN_RATE_LIMIT_PER_MIN", 10),
		SeedOwnerEmail:       getEnv("SEED_OWNER_EMAIL", "owner@icecream.local"),
		SeedOwnerPassword:    getEnv("SEED_OWNER_PASSWORD", "Owner@12345"),
		SeedManagerEmail:     getEnv("SEED_MANAGER_EMAIL", "manager@icecream.local"),
		SeedManagerPassword:  getEnv("SEED_MANAGER_PASSWORD", "Manager@12345"),
		SeedCashierEmail:     getEnv("SEED_CASHIER_EMAIL", "cashier@icecream.local"),
		SeedCashierPassword:  getEnv("SEED_CASHIER_PASSWORD", "Cashier@12345"),
	}

	if cfg.AppEnv != "test" {
		if strings.TrimSpace(cfg.JWTSecret) == "" {
			return nil, errors.New("JWT_SECRET is required")
		}
		if len(cfg.JWTSecret) < 32 {
			return nil, errors.New("JWT_SECRET must be at least 32 characters")
		}
		if strings.TrimSpace(cfg.DatabaseURL) == "" {
			return nil, errors.New("DATABASE_URL is required")
		}
	} else {
		if strings.TrimSpace(cfg.JWTSecret) == "" {
			cfg.JWTSecret = "test-secret-key-that-is-at-least-32-chars!"
		}
		if len(cfg.JWTSecret) < 32 {
			cfg.JWTSecret = "test-secret-key-that-is-at-least-32-chars!"
		}
	}

	if cfg.AppPort == "" {
		return nil, errors.New("APP_PORT is required")
	}
	if _, err := strconv.Atoi(cfg.AppPort); err != nil {
		return nil, fmt.Errorf("APP_PORT must be a valid port: %w", err)
	}

	if cfg.BcryptCost < 4 || cfg.BcryptCost > 31 {
		cfg.BcryptCost = 12
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		v = strings.ToLower(v)
		if v == "true" || v == "1" || v == "yes" {
			return true
		}
		if v == "false" || v == "0" || v == "no" {
			return false
		}
	}
	return def
}
