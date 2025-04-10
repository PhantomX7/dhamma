package utility

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Common repository errors
var (
	ErrNotFound    = errors.New("record not found")
	ErrDuplicate   = errors.New("duplicate record")
	ErrInvalidData = errors.New("invalid data")
	ErrForeignKey  = errors.New("foreign key constraint failed")
	ErrDatabase    = errors.New("database error")
)

func LogError(errString string, err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Print(errString, ": ", err)
	}
	return err
}

// WrapError wraps an error with additional context
func WrapError(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)
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
