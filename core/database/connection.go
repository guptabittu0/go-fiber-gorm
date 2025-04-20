package database

import (
	"fmt"
	"go-fiber-gorm/config"
	"go-fiber-gorm/core/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Connection is the database connection manager
type Connection struct {
	DB *gorm.DB
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
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

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL database -> Failed to Connect \n\t %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("PostgresSQL database -> Failed to Connect \n\t %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	logger.Info("PostgreSQL database -> Connected")

	return &Connection{DB: db}, nil
}

// GetDB returns the database instance
func (c *Connection) GetDB() *gorm.DB {
	return c.DB
}

// Close closes the database connection
func (c *Connection) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate runs auto migrations for the provided models
func (c *Connection) AutoMigrate(models ...interface{}) error {
	logger.Info("Database migrations -> Running...")
	if err := c.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("database migration -> %w", err)
	}
	logger.Info("Database migrations -> Completed")
	return nil
}
