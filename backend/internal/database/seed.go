package database

import (
	"strings"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/MrVirul/icecream-inventory/internal/models"
	"gorm.io/gorm"
)

type SeedConfig struct {
	AppEnv              string
	BcryptCost          int
	SeedOwnerEmail      string
	SeedOwnerPassword   string
	SeedManagerEmail    string
	SeedManagerPassword string
	SeedCashierEmail    string
	SeedCashierPassword string
}

func Seed(db *gorm.DB, cfg SeedConfig) error {
	// Seed users always
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Owner", strings.ToLower(strings.TrimSpace(cfg.SeedOwnerEmail)), cfg.SeedOwnerPassword, "Owner"},
		{"Manager", strings.ToLower(strings.TrimSpace(cfg.SeedManagerEmail)), cfg.SeedManagerPassword, "Manager"},
		{"Cashier", strings.ToLower(strings.TrimSpace(cfg.SeedCashierEmail)), cfg.SeedCashierPassword, "Cashier"},
	}
	for _, u := range users {
		var existing models.User
		err := db.Where("LOWER(email) = LOWER(?)", u.Email).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			hash, err := auth.HashPassword(u.Password, cfg.BcryptCost)
			if err != nil {
				return err
			}
			user := models.User{
				Name:         u.Name,
				Email:        u.Email,
				PasswordHash: hash,
				Role:         u.Role,
			}
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		// if exists, do nothing (idempotent)
	}

	// Seed sample items and customers only if not production
	if cfg.AppEnv == "production" {
		return nil
	}

	// Sample items
	items := []models.InventoryItem{
		{Name: "Chocolate", Category: "flavor", Description: "Rich dark chocolate", Quantity: 12, Unit: "tub", PriceCents: 1850, ReorderLevel: 5},
		{Name: "Vanilla", Category: "flavor", Description: "Classic vanilla", Quantity: 20, Unit: "tub", PriceCents: 1650, ReorderLevel: 5},
		{Name: "Strawberry", Category: "flavor", Description: "Fresh strawberry", Quantity: 3, Unit: "tub", PriceCents: 1750, ReorderLevel: 5},
		{Name: "Mint Chip", Category: "flavor", Description: "Mint with chocolate chips", Quantity: 8, Unit: "tub", PriceCents: 1850, ReorderLevel: 5},
		{Name: "Butterscotch", Category: "flavor", Description: "Butterscotch flavor", Quantity: 10, Unit: "tub", PriceCents: 1750, ReorderLevel: 5},
		{Name: "Sprinkles", Category: "topping", Description: "Rainbow sprinkles", Quantity: 50, Unit: "box", PriceCents: 250, ReorderLevel: 10},
		{Name: "Waffle Cone", Category: "cone", Description: "Crispy waffle cone", Quantity: 100, Unit: "pcs", PriceCents: 150, ReorderLevel: 20},
		{Name: "Sugar Cone", Category: "cone", Description: "Sweet sugar cone", Quantity: 80, Unit: "pcs", PriceCents: 100, ReorderLevel: 20},
		{Name: "Paper Cup", Category: "cup", Description: "Paper cup for scoops", Quantity: 200, Unit: "pcs", PriceCents: 50, ReorderLevel: 30},
		{Name: "Cone Sleeve", Category: "packaging", Description: "Sleeve for cones", Quantity: 150, Unit: "pcs", PriceCents: 30, ReorderLevel: 30},
	}
	for _, it := range items {
		var existing models.InventoryItem
		if err := db.Where("LOWER(name) = LOWER(?)", it.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&it).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	// Sample customers
	custs := []models.Customer{
		{Name: "John Doe", Phone: "1234567890", Address: "123 Main St"},
		{Name: "Jane Smith", Phone: "9876543210", Address: "456 Oak Ave"},
		{Name: "Bob Johnson", Phone: "5551234567", Address: "789 Pine Rd"},
		{Name: "Alice Brown", Phone: "5559876543", Address: "321 Elm St"},
		{Name: "Charlie Wilson", Phone: "5554567890", Address: "654 Maple Dr"},
	}
	for _, c := range custs {
		var existing models.Customer
		if err := db.Where("phone = ?", c.Phone).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&c).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}
