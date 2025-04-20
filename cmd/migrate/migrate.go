package main

import (
	"flag"
	"go-fiber-gorm/config"
	"go-fiber-gorm/core/database"
	"go-fiber-gorm/core/logger"
	"go-fiber-gorm/migrations"
	"os"
)

func main() {
	// Define command line flags
	var down bool
	flag.BoolVar(&down, "down", false, "Roll back the last migration")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	dbConn, err := database.NewConnection(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	if down {
		// Rollback last migration
		logger.Info("Rolling back last migration...")
		if err := migrations.RollbackLastMigration(dbConn.DB); err != nil {
			logger.Fatal("Failed to rollback migration:", err)
		}
		logger.Info("Migration rollback completed successfully")
	} else {
		// Run migrations
		logger.Info("Running database migrations...")
		if err := migrations.RunMigrations(dbConn.DB); err != nil {
			logger.Fatal("Failed to run migrations:", err)
		}
		logger.Info("Migrations completed successfully")
	}

	os.Exit(0)
}
