package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/MrVirul/icecream-inventory/internal/config"
	"github.com/MrVirul/icecream-inventory/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}
	db, err := database.Open(cfg.DatabaseURL, cfg.AppEnv)
	if err != nil {
		slog.Error("db open failed", "error", err)
		os.Exit(1)
	}
	if err := database.Seed(db, database.SeedConfig{
		AppEnv:              cfg.AppEnv,
		BcryptCost:          cfg.BcryptCost,
		SeedOwnerEmail:      cfg.SeedOwnerEmail,
		SeedOwnerPassword:   cfg.SeedOwnerPassword,
		SeedManagerEmail:    cfg.SeedManagerEmail,
		SeedManagerPassword: cfg.SeedManagerPassword,
		SeedCashierEmail:    cfg.SeedCashierEmail,
		SeedCashierPassword: cfg.SeedCashierPassword,
	}); err != nil {
		slog.Error("seed failed", "error", err)
		os.Exit(1)
	}
	fmt.Println("seeding completed")
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
