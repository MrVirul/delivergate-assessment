package database

import (
	"errors"
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

type seedAccount struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func Seed(db *gorm.DB, cfg SeedConfig) error {
	// Seed users always
	users := []seedAccount{
		{"Owner", normalizeEmail(cfg.SeedOwnerEmail), cfg.SeedOwnerPassword, "Owner"},
		{"Manager", normalizeEmail(cfg.SeedManagerEmail), cfg.SeedManagerPassword, "Manager"},
		{"Cashier", normalizeEmail(cfg.SeedCashierEmail), cfg.SeedCashierPassword, "Cashier"},
	}
	for _, u := range users {
		if err := seedUser(db, u, cfg.BcryptCost); err != nil {
			return err
		}
	}

	// Seed sample items and customers only if not production
	if cfg.AppEnv == "production" {
		return nil
	}

	if err := seedItems(db); err != nil {
		return err
	}

	return seedCustomers(db)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func seedUser(db *gorm.DB, u seedAccount, bcryptCost int) error {
	var existing models.User
	err := db.Where("LOWER(email) = LOWER(?)", u.Email).First(&existing).Error
	if err == nil {
		// if exists, do nothing (idempotent)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := auth.HashPassword(u.Password, bcryptCost)
	if err != nil {
		return err
	}
	return db.Create(&models.User{
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: hash,
		Role:         u.Role,
	}).Error
}

func seedItems(db *gorm.DB) error {
	items := []models.InventoryItem{
		{Name: "Ceylon Vanilla", Category: "flavor", Description: "Classic ice cream made with Ceylon vanilla", Quantity: 20, Unit: "tub", PriceCents: 145000, ReorderLevel: 5},
		{Name: "Kithul Treacle", Category: "flavor", Description: "Sri Lankan kithul treacle flavoured ice cream", Quantity: 14, Unit: "tub", PriceCents: 165000, ReorderLevel: 5},
		{Name: "Wood Apple", Category: "flavor", Description: "Tangy wood apple (divul) ice cream", Quantity: 9, Unit: "tub", PriceCents: 155000, ReorderLevel: 5},
		{Name: "King Coconut", Category: "flavor", Description: "Refreshing king coconut (thambili) ice cream", Quantity: 12, Unit: "tub", PriceCents: 150000, ReorderLevel: 5},
		{Name: "Mango", Category: "flavor", Description: "Sweet Sri Lankan mango ice cream", Quantity: 3, Unit: "tub", PriceCents: 140000, ReorderLevel: 5},
		{Name: "Curd and Honey", Category: "flavor", Description: "Buffalo curd with kithul honey", Quantity: 8, Unit: "tub", PriceCents: 175000, ReorderLevel: 5},
		{Name: "Kithul Treacle Syrup", Category: "topping", Description: "Pure kithul treacle for drizzling", Quantity: 40, Unit: "box", PriceCents: 85000, ReorderLevel: 10},
		{Name: "Waffle Cone", Category: "cone", Description: "Crispy waffle cone", Quantity: 100, Unit: "pcs", PriceCents: 4500, ReorderLevel: 20},
		{Name: "Paper Cup", Category: "cup", Description: "Paper cup for scoops", Quantity: 200, Unit: "pcs", PriceCents: 1500, ReorderLevel: 30},
		{Name: "Banana Leaf Wrap", Category: "packaging", Description: "Eco-friendly banana leaf wrapping", Quantity: 150, Unit: "pcs", PriceCents: 1000, ReorderLevel: 30},
	}
	for _, it := range items {
		var existing models.InventoryItem
		err := db.Where("LOWER(name) = LOWER(?)", it.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&it).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedCustomers(db *gorm.DB) error {
	customers := []models.Customer{
		{Name: "Nimal Perera", Email: stringPtr("nimal.perera@gmail.com"), Phone: "+94 77 123 4567", Address: "No. 12, Galle Road, Colombo 03"},
		{Name: "Sanduni Fernando", Email: stringPtr("sanduni.fernando@gmail.com"), Phone: "+94 71 234 5678", Address: "No. 45, Peradeniya Road, Kandy"},
		{Name: "Kasun Silva", Email: stringPtr("kasun.silva@gmail.com"), Phone: "+94 76 345 6789", Address: "No. 8, Beach Road, Galle"},
		{Name: "Dilini Jayawardena", Email: stringPtr("dilini.jayawardena@gmail.com"), Phone: "+94 72 456 7890", Address: "No. 23, Temple Road, Negombo"},
		{Name: "Chamara Bandara", Email: stringPtr("chamara.bandara@gmail.com"), Phone: "+94 75 567 8901", Address: "No. 67, Hospital Road, Jaffna"},
	}
	for _, c := range customers {
		var existing models.Customer
		err := db.Where("phone = ?", c.Phone).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

func stringPtr(v string) *string {
	return &v
}
