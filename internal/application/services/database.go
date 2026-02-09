package services

import "context"

// Transactor defines an interface for running operations within a database transaction.
type Transactor interface {
	// RunInTransaction executes the given function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function returns nil, the transaction is committed.
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
