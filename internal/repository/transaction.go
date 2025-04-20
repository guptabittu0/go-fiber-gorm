package repository

import (
	"context"
	"go-fiber-gorm/pkg/logger"

	"gorm.io/gorm"
)

// TxManager is responsible for managing database transactions
type TxManager struct {
	db *gorm.DB
}

// NewTxManager creates a new transaction manager
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

// WithTransaction executes the given function within a transaction
// If the function returns an error, the transaction is rolled back
// If the function returns nil, the transaction is committed
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// Begin transaction
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Error("Failed to begin transaction:", tx.Error)
		return tx.Error
	}

	// Execute the function
	err := fn(tx)
	if err != nil {
		// Rollback transaction in case of error
		logger.Info("Rolling back transaction due to error:", err)
		if rbErr := tx.Rollback().Error; rbErr != nil {
			logger.Error("Failed to rollback transaction:", rbErr)
			// Return original error, not rollback error
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction:", err)
		return err
	}

	return nil
}
