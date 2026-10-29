package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-that-is-at-least-32-chars!")
	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "test", cfg.AppEnv)
	require.Equal(t, 60, cfg.JWTTTLMinutes)
}

func TestLoadMissingJWTInNonTest(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_ENV", "development")
	os.Setenv("DATABASE_URL", "postgres://x")
	_, err := Load()
	require.Error(t, err)
}

func TestLoadShortJWT(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_ENV", "development")
	os.Setenv("DATABASE_URL", "postgres://x")
	os.Setenv("JWT_SECRET", "short")
	_, err := Load()
	require.Error(t, err)
}
