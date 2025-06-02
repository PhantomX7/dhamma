package errors

import (
	"fmt"
	"net/http"
)

// ErrorType defines the type of error for proper handling
type ErrorType string

// Error type constants
const (
	ErrorTypeValidation         ErrorType = "VALIDATION_ERROR"
	ErrorTypeAuthentication     ErrorType = "AUTHENTICATION_ERROR"
	ErrorTypeAuthorization      ErrorType = "AUTHORIZATION_ERROR"
	ErrorTypeNotFound           ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict           ErrorType = "CONFLICT_ERROR"
	ErrorTypeDatabase           ErrorType = "DATABASE_ERROR"
	ErrorTypeExternal           ErrorType = "EXTERNAL_SERVICE_ERROR"
	ErrorTypeInternal           ErrorType = "INTERNAL_ERROR"
	ErrorTypeRateLimit          ErrorType = "RATE_LIMIT_ERROR"
	ErrorTypeBusiness           ErrorType = "BUSINESS_LOGIC_ERROR"
	ErrorTypeServiceUnavailable ErrorType = "SERVICE_UNAVAILABLE_ERROR"
	ErrorTypeTooManyRequests    ErrorType = "TOO_MANY_REQUESTS_ERROR"
)

// Error code constants
const (
	CodeValidationFailed        = "E1001"
	CodeInvalidCredentials      = "E1002"
	CodeTokenExpired            = "E1003"
	CodeInsufficientPermissions = "E1004"
	CodeResourceNotFound        = "E1005"
	CodeResourceAlreadyExists   = "E1006"
	CodeDatabaseConnection      = "E1007"
	CodeDatabaseQuery           = "E1008"
	CodeExternalServiceDown     = "E1009"
	CodeExternalServiceTimeout  = "E1010"
	CodeInternalServerError     = "E1011"
	CodeRateLimitExceeded       = "E1012"
	CodeBusinessRuleViolation   = "E1013"
	CodeServiceUnavailable      = "E1014"
	CodeTooManyRequests         = "E1015"
)

// AppError represents an application error with enhanced context
type AppError struct {
	Message   string                 `json:"message"`
	Code      string                 `json:"code"`
	Type      ErrorType              `json:"type"`
	Err       error                  `json:"-"` // Don't expose internal error in JSON
	Status    int                    `json:"-"` // Don't expose HTTP status in JSON
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s [%s]: %s - %v", e.Type, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s [%s]: %s", e.Type, e.Code, e.Message)
}

// WithRequestID adds request ID to the error
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// Validation Errors
func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeValidationFailed,
		Type:    ErrorTypeValidation,
		Err:     err,
		Status:  http.StatusBadRequest,
	}
}

// Authentication Errors
func NewAuthenticationError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeInvalidCredentials,
		Type:    ErrorTypeAuthentication,
		Err:     err,
		Status:  http.StatusUnauthorized,
	}
}

func NewTokenExpiredError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeTokenExpired,
		Type:    ErrorTypeAuthentication,
		Err:     err,
		Status:  http.StatusUnauthorized,
	}
}

// Authorization Errors
func NewAuthorizationError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeInsufficientPermissions,
		Type:    ErrorTypeAuthorization,
		Err:     err,
		Status:  http.StatusForbidden,
	}
}

// Not Found Errors
func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeResourceNotFound,
		Type:    ErrorTypeNotFound,
		Err:     err,
		Status:  http.StatusNotFound,
	}
}

// Conflict Errors
func NewConflictError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeResourceAlreadyExists,
		Type:    ErrorTypeConflict,
		Err:     err,
		Status:  http.StatusConflict,
	}
}

// Database Errors
func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeDatabaseQuery,
		Type:    ErrorTypeDatabase,
		Err:     err,
		Status:  http.StatusInternalServerError,
	}
}

// External Service Errors
func NewExternalServiceError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeExternalServiceDown,
		Type:    ErrorTypeExternal,
		Err:     err,
		Status:  http.StatusBadGateway,
	}
}

// Rate Limit Errors
func NewRateLimitError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeRateLimitExceeded,
		Type:    ErrorTypeRateLimit,
		Err:     err,
		Status:  http.StatusTooManyRequests,
	}
}

// Business Logic Errors
func NewBusinessError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeBusinessRuleViolation,
		Type:    ErrorTypeBusiness,
		Err:     err,
		Status:  http.StatusUnprocessableEntity,
	}
}

// Internal Server Errors
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeInternalServerError,
		Type:    ErrorTypeInternal,
		Err:     err,
		Status:  http.StatusInternalServerError,
	}
}

// Service Unavailable Errors
func NewServiceUnavailableError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeServiceUnavailable,
		Type:    ErrorTypeServiceUnavailable,
		Err:     err,
		Status:  http.StatusServiceUnavailable,
	}
}

// Too Many Requests Errors
func NewTooManyRequestsError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    CodeTooManyRequests,
		Type:    ErrorTypeTooManyRequests,
		Err:     err,
		Status:  http.StatusTooManyRequests,
	}
}

// Legacy functions for backward compatibility
func NewServiceError(message string, err error) *AppError {
	return NewBusinessError(message, err)
}

func NewBadRequestError(message string, err error) *AppError {
	return NewValidationError(message, err)
}
