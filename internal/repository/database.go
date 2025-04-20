package repository

import (
	"fmt"
	"go-fiber-gorm/config"
	"go-fiber-gorm/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Database is the database connection
var DB *gorm.DB

// ConnectDatabase establishes a connection to the database
func ConnectDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	logger.Info("PostgreSQL database -> Connecting...")

	// Build the DSN
	dsn := cfg.GetDSN()

	// Connect to the database
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "app_", // Table name prefix
			SingularTable: false,  // Use plural form for table names
		},
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL database -> Failed to Connect \n\t %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("PostgresSQL database -> Failed to Connect \n\t %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	logger.Info("PostgreSQL database -> Connected")
	return DB, nil
}

// AutoMigrate runs auto migrations for the provided models
func AutoMigrate(models ...interface{}) error {
	logger.Info("Database migrations -> Running...")
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("database migration -> %w", err)
	}
	logger.Info("Database migrations -> Completed")
	return nil
}
