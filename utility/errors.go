package utility

import (
	"errors"
	"fmt"

	"github.com/PhantomX7/dhamma/utility/logger"
	"go.uber.org/zap"
)

// Common repository errors
var (
	ErrNotFound         = errors.New("record not found")
	ErrDuplicate        = errors.New("duplicate record")
	ErrInvalidData      = errors.New("invalid data")
	ErrForeignKey       = errors.New("foreign key constraint failed")
	ErrDatabase         = errors.New("database error")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidRequest   = errors.New("invalid request")
)

// WrapError wraps an error with additional context
func WrapError(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)
	logger.Get().Error(message,
		zap.Error(err),
	)
	return fmt.Errorf("%s: %w", message, err)
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsDuplicate checks if the error is a duplicate error
func IsDuplicate(err error) bool {
	return errors.Is(err, ErrDuplicate)
}
