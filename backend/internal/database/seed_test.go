package database

import (
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/MrVirul/icecream-inventory/internal/models"
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
		SeedOwnerEmail:      "owner@mail.com",
		SeedOwnerPassword:   "Owner@123",
		SeedManagerEmail:    "manager@mail.com",
		SeedManagerPassword: "Manager@123",
		SeedCashierEmail:    "cashier@mail.com",
		SeedCashierPassword: "Cashier@123",
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
	require.NoError(t, db.Table("users").Select("password_hash").Where("email = ?", "owner@mail.com").First(&u).Error)
	require.NoError(t, auth.ComparePassword(u.PasswordHash, "Owner@123"))

	// verify sample items and customers are seeded once
	var itemCount, customerCount int64
	db.Model(&models.InventoryItem{}).Count(&itemCount)
	db.Model(&models.Customer{}).Count(&customerCount)
	require.Equal(t, int64(10), itemCount)
	require.Equal(t, int64(5), customerCount)

	var lowStock models.InventoryItem
	require.NoError(t, db.Where("name = ?", "Mango").First(&lowStock).Error)
	require.Less(t, lowStock.Quantity, lowStock.ReorderLevel)

	var customer models.Customer
	require.NoError(t, db.Where("phone = ?", "+94 77 123 4567").First(&customer).Error)
	require.NotNil(t, customer.Email)
	require.Contains(t, customer.Address, "Colombo")
}
