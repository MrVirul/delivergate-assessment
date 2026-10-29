package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MrVirul/icecream-inventory/internal/config"
	"github.com/MrVirul/icecream-inventory/internal/database"
	"github.com/MrVirul/icecream-inventory/internal/router"
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
	if cfg.AutoMigrate {
		if err := database.Migrate(db); err != nil {
			slog.Error("migrate failed", "error", err)
			os.Exit(1)
		}
	}

	app := router.New(cfg, db)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.AppPort)
		slog.Info("starting server", "addr", addr)
		if err := app.Listen(addr); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}
