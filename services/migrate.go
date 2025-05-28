package services

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"

	"github.com/mrgThang/flashcard-be/config"
)

func RunMigrations(migrationsDir string) {
	cfg := &config.Config{}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dsn := cfg.MysqlConfig.DSN()

	sourceURL := fmt.Sprintf("file://%s", migrationsDir)
	m, err := migrate.New(
		sourceURL,
		"mysql://"+dsn,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations applied successfully.")
}
