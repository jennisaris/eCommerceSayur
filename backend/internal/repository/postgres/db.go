package postgres

import (
	"ecommerce-backend/internal/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	logLevel := logger.Warn
	if os.Getenv("APP_ENV") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// High concurrency tuning for k6 1,000+ VU testing
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

// RunMigrations executes up SQL migration files from the given directory
func RunMigrations(db *gorm.DB, migrationsDir string, runSeeder bool) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		// Log and continue if dir not found (e.g. running in testing container where initdb ran it)
		log.Printf("Migration dir read info: %v", err)
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" && len(entry.Name()) > 7 && entry.Name()[len(entry.Name())-7:] == ".up.sql" {
			// In production, skip seeder files unless explicitly enabled
			if !runSeeder && strings.Contains(strings.ToLower(entry.Name()), "seed") {
				log.Printf("ℹ️ Skipping seed data migration in production mode: %s", entry.Name())
				continue
			}

			filePath := filepath.Join(migrationsDir, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed reading migration %s: %w", entry.Name(), err)
			}

			if err := db.Exec(string(content)).Error; err != nil {
				log.Printf("Notice on executing migration %s: %v", entry.Name(), err)
			} else {
				log.Printf("Successfully applied migration: %s", entry.Name())
			}
		}
	}

	return nil
}
