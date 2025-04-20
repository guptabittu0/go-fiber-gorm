package migrations

import (
	"go-fiber-gorm/core/logger"
	"go-fiber-gorm/modules/auth"
	"go-fiber-gorm/modules/user"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	Name     string
	Migrate  func(*gorm.DB) error
	Rollback func(*gorm.DB) error
}

// Migrations is a list of all migrations
var Migrations = []Migration{
	{
		Name: "create_users_table",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&user.User{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&user.User{})
		},
	},
	{
		Name: "create_sessions_table",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&auth.Session{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&auth.Session{})
		},
	},
	// Add more migrations as needed
}

// RunMigrations runs all migrations
func RunMigrations(db *gorm.DB) error {
	// Create migrations table if it doesn't exist
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return err
	}

	// Get executed migrations
	var executed []MigrationRecord
	if err := db.Find(&executed).Error; err != nil {
		return err
	}

	// Convert to a map for easier lookup
	executedMap := make(map[string]bool)
	for _, m := range executed {
		executedMap[m.Name] = true
	}

	// Run pending migrations
	for _, migration := range Migrations {
		if !executedMap[migration.Name] {
			logger.Info("Running migration:", migration.Name)

			// Start a transaction
			tx := db.Begin()

			// Run the migration
			if err := migration.Migrate(tx); err != nil {
				tx.Rollback()
				return err
			}

			// Record the migration
			if err := tx.Create(&MigrationRecord{Name: migration.Name}).Error; err != nil {
				tx.Rollback()
				return err
			}

			// Commit the transaction
			if err := tx.Commit().Error; err != nil {
				return err
			}

			logger.Info("Migration completed:", migration.Name)
		}
	}

	return nil
}

// RollbackLastMigration rolls back the last migration
func RollbackLastMigration(db *gorm.DB) error {
	// Get last executed migration
	var lastMigration MigrationRecord
	if err := db.Order("id desc").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("No migrations to rollback")
			return nil
		}
		return err
	}

	// Find the migration in our list
	var migration Migration
	found := false
	for _, m := range Migrations {
		if m.Name == lastMigration.Name {
			migration = m
			found = true
			break
		}
	}

	if !found {
		logger.Warn("Migration not found in list:", lastMigration.Name)
		return nil
	}

	// Start a transaction
	tx := db.Begin()

	// Run the rollback
	logger.Info("Rolling back migration:", migration.Name)
	if err := migration.Rollback(tx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the migration record
	if err := tx.Delete(&lastMigration).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	logger.Info("Rollback completed for:", migration.Name)
	return nil
}

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}
