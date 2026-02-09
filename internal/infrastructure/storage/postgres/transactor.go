package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "gorm_tx"

var ErrNoTransaction = errors.New("no transaction in context")

// Transactor represents a transaction manager for GORM
type Transactor struct {
	db *gorm.DB
}

// NewTransactor creates a new Transactor instance
func NewTransactor(db *gorm.DB) *Transactor {
	return &Transactor{db: db}
}

// RunInTransaction executes the given function within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
func (t *Transactor) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx)
	})
}

// GetTx retrieves the transaction from the context
// Returns the transaction if exists, otherwise returns the original DB instance
func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db
}

// EnsureTx checks that a transaction exists in the context
// Returns an error if no transaction is found
func EnsureTx(ctx context.Context) error {
	if _, ok := ctx.Value(txKey).(*gorm.DB); !ok {
		return ErrNoTransaction
	}
	return nil
}
