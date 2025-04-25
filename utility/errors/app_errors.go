package errors

import "net/http"

// ErrorType defines the type of error for proper handling
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeService represents service-level errors
	ErrorTypeService ErrorType = "service"
	// ErrorTypeBadRequest represents general bad request errors
	ErrorTypeBadRequest ErrorType = "bad_request"
)

// AppError represents an application error with context
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
	Status  int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// NewValidationError creates a validation error
func NewValidationError(err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: "Validation failed",
		Err:     err,
		Status:  http.StatusBadRequest,
	}
}

// NewServiceError creates a service error
func NewServiceError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeService,
		Message: message,
		Err:     err,
		Status:  http.StatusUnprocessableEntity,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Err:     err,
		Status:  http.StatusBadRequest,
	}
}
