package repositories

import "errors"

var (
	// ErrNotFound indicates that the requested entity was not found in the repository
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidArgument indicates that an invalid argument was provided to a repository method
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrConnectionFailed indicates that the database connection is broken or unavailable
	ErrConnectionFailed = errors.New("database connection failed")

	// ErrDuplicateKey indicates that a unique constraint was violated
	ErrDuplicateKey = errors.New("duplicate key constraint violation")

	// ErrForeignKeyViolation indicates that a foreign key constraint was violated
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")

	// ErrCheckConstraintViolation indicates that a check constraint was violated
	ErrCheckConstraintViolation = errors.New("check constraint violation")

	// ErrOperationFailed indicates that a repository operation failed for an unknown reason
	ErrOperationFailed = errors.New("repository operation failed")
)
