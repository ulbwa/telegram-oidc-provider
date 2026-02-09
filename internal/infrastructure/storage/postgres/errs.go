package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

// isConnectionError checks if the error is related to database connection issues
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for context errors (timeout, cancellation)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	
	// Check for driver.ErrBadConn
	if errors.Is(err, driver.ErrBadConn) {
		return true
	}
	
	// Check for GORM invalid DB error (could indicate connection issues)
	if errors.Is(err, gorm.ErrInvalidDB) {
		return true
	}
	
	// Check for common connection error messages as fallback
	errMsg := strings.ToLower(err.Error())
	connectionErrors := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"no connection",
		"connection closed",
		"connection lost",
		"dial tcp",
		"i/o timeout",
		"network is unreachable",
		"sql: database is closed",
	}
	
	for _, connErr := range connectionErrors {
		if strings.Contains(errMsg, connErr) {
			return true
		}
	}
	
	return false
}

// mapError maps database errors to repository errors
func mapError(err error, operation string) error {
	if err == nil {
		return nil
	}
	
	// Check GORM-specific errors first
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: %s: %v", repositories.ErrNotFound, operation, err)
	}
	
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: %s: %v", repositories.ErrDuplicateKey, operation, err)
	}
	
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return fmt.Errorf("%w: %s: %v", repositories.ErrForeignKeyViolation, operation, err)
	}
	
	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		return fmt.Errorf("%w: %s: %v", repositories.ErrCheckConstraintViolation, operation, err)
	}
	
	// Check for connection-related errors
	if isConnectionError(err) {
		return fmt.Errorf("%w: %s: %v", repositories.ErrConnectionFailed, operation, err)
	}
	
	// Default to operation failed for all other errors
	return fmt.Errorf("%w: %s: %v", repositories.ErrOperationFailed, operation, err)
}
