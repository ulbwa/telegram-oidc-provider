package service

import "context"

// Transactor defines transaction boundary for application use-cases.
type Transactor interface {
	// RunInTransaction executes the given function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function returns nil, the transaction is committed.
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
