package repository

import "errors"

// Repository errors
var (
	// ErrDatabaseError is returned for general database errors
	ErrDatabaseError = errors.New("database error")

	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("entity not found")

	// ErrDuplicate is returned when trying to create/update an entity with a duplicate unique field
	ErrDuplicate = errors.New("duplicate entity")

	// ErrCorruptedData is returned when data retrieved from database is corrupted or invalid
	ErrCorruptedData = errors.New("corrupted data in database")

	// ErrEncryptionFailed is returned when encryption or decryption fails
	ErrEncryptionFailed = errors.New("encryption operation failed")

	// ErrNoTransaction is returned when no transaction exists in context
	ErrNoTransaction = errors.New("no transaction in context")
)
