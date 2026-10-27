package database

import (
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSeedIdempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, Migrate(db))
	cfg := SeedConfig{
		AppEnv:              "development",
		BcryptCost:          4,
		SeedOwnerEmail:      "owner@icecream.local",
		SeedOwnerPassword:   "Owner@12345",
		SeedManagerEmail:    "manager@icecream.local",
		SeedManagerPassword: "Manager@12345",
		SeedCashierEmail:    "cashier@icecream.local",
		SeedCashierPassword: "Cashier@12345",
	}
	require.NoError(t, Seed(db, cfg))
	require.NoError(t, Seed(db, cfg))
	var count int64
	db.Model(&struct{}{}).Table("users").Count(&count)
	require.Equal(t, int64(3), count)
	// verify password
	var u struct {
		PasswordHash string
	}
	require.NoError(t, db.Table("users").Select("password_hash").Where("email = ?", "owner@icecream.local").First(&u).Error)
	require.NoError(t, auth.ComparePassword(u.PasswordHash, "Owner@12345"))
}
