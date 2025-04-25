package errors

import "net/http"

// ErrorType defines the type of error for proper handling
type ErrorType string

// AppError represents an application error with context
type AppError struct {
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

// NewServiceError creates a service error
func NewServiceError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Err:     err,
		Status:  http.StatusUnprocessableEntity,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Err:     err,
		Status:  http.StatusBadRequest,
	}
}
