package database

import (
	"github.com/MrVirul/icecream-inventory/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.InventoryItem{},
		&models.Customer{},
		&models.RevokedToken{},
	)
}
