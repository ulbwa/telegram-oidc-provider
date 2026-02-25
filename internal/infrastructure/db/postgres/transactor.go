package postgres

import (
	"context"

	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "gorm_tx"

// GormTransactor is a GORM-based implementation of service.Transactor.
type GormTransactor struct {
	db *gorm.DB
}

// Compile-time check that GormTransactor implements service.Transactor
var _ service.Transactor = (*GormTransactor)(nil)

// NewGormTransactor creates a new GORM transaction manager.
func NewGormTransactor(db *gorm.DB) *GormTransactor {
	return &GormTransactor{db: db}
}

// RunInTransaction executes the given function within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
func (t *GormTransactor) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return fn(ctx)
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx)
	})
}

// GetTx retrieves the transaction from the context.
// Returns the transaction if it exists, otherwise returns the original DB instance.
func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db
}

// EnsureTx checks that a transaction exists in the context.
// Returns an error if no transaction is found.
func EnsureTx(ctx context.Context) error {
	if _, ok := ctx.Value(txKey).(*gorm.DB); !ok {
		return repository.ErrNoTransaction
	}
	return nil
}
