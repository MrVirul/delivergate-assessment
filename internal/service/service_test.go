package service

import (
	"testing"

	"github.com/MrVirul/icecream-inventory/internal/database"
	"github.com/MrVirul/icecream-inventory/internal/dto"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/MrVirul/icecream-inventory/internal/repository"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.Migrate(db))
	return db
}

func TestItemServiceCreateDuplicate(t *testing.T) {
	db := setupTestDB(t)
	svc := NewItemService(repository.NewItemRepository(db))
	_, err := svc.Create(dto.ItemCreateRequest{Name: "Chocolate", Category: "flavor", Quantity: intp(5), Unit: "tub", PriceCents: int64p(1000)}, 1)
	require.NoError(t, err)
	_, err = svc.Create(dto.ItemCreateRequest{Name: "chocolate", Category: "flavor", Quantity: intp(5), Unit: "tub", PriceCents: int64p(1000)}, 1)
	var ae *httpx.AppError
	require.ErrorAs(t, err, &ae)
	require.Equal(t, 409, ae.Status)
}

func TestItemLowStock(t *testing.T) {
	db := setupTestDB(t)
	svc := NewItemService(repository.NewItemRepository(db))
	item, err := svc.Create(dto.ItemCreateRequest{Name: "X", Category: "flavor", Quantity: intp(2), Unit: "tub", PriceCents: int64p(1000), ReorderLevel: intp(5)}, 1)
	require.NoError(t, err)
	require.True(t, item.LowStock)
}

func intp(i int) *int       { return &i }
func int64p(i int64) *int64 { return &i }

func TestCustomerEmailLowercase(t *testing.T) {
	db := setupTestDB(t)
	svc := NewCustomerService(repository.NewCustomerRepository(db))
	e := "Test@Example.COM"
	c, err := svc.Create(dto.CustomerCreateRequest{Name: "Alice", Email: &e, Phone: "1234567"}, 1)
	require.NoError(t, err)
	require.NotNil(t, c.Email)
	require.Equal(t, "test@example.com", *c.Email)
}
