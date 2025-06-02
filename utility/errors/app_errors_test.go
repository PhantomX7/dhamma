package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "error with underlying error",
			appError: &AppError{
				Message: "validation failed",
				Code:    CodeValidationFailed,
				Type:    ErrorTypeValidation,
				Err:     errors.New("field is required"),
			},
			expected: "VALIDATION_ERROR [E1001]: validation failed - field is required",
		},
		{
			name: "error without underlying error",
			appError: &AppError{
				Message: "user not found",
				Code:    CodeResourceNotFound,
				Type:    ErrorTypeNotFound,
				Err:     nil,
			},
			expected: "NOT_FOUND_ERROR [E1005]: user not found",
		},
		{
			name: "empty message",
			appError: &AppError{
				Message: "",
				Code:    CodeInternalServerError,
				Type:    ErrorTypeInternal,
				Err:     nil,
			},
			expected: "INTERNAL_ERROR [E1011]: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAppError_WithRequestID(t *testing.T) {
	appError := &AppError{
		Message: "test error",
		Code:    CodeValidationFailed,
		Type:    ErrorTypeValidation,
		Status:  http.StatusBadRequest,
	}

	requestID := "req-123456"
	result := appError.WithRequestID(requestID)

	// Should return the same instance
	assert.Same(t, appError, result)
	assert.Equal(t, requestID, result.RequestID)
}

func TestAppError_WithDetails(t *testing.T) {
	appError := &AppError{
		Message: "test error",
		Code:    CodeValidationFailed,
		Type:    ErrorTypeValidation,
		Status:  http.StatusBadRequest,
	}

	details := map[string]interface{}{
		"field": "email",
		"value": "invalid-email",
	}

	result := appError.WithDetails(details)

	// Should return the same instance
	assert.Same(t, appError, result)
	assert.Equal(t, details, result.Details)
}

func TestNewValidationError(t *testing.T) {
	message := "validation failed"
	underlyingErr := errors.New("field is required")

	result := NewValidationError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeValidationFailed, result.Code)
	assert.Equal(t, ErrorTypeValidation, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusBadRequest, result.Status)
}

func TestNewAuthenticationError(t *testing.T) {
	message := "invalid credentials"
	underlyingErr := errors.New("password mismatch")

	result := NewAuthenticationError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeInvalidCredentials, result.Code)
	assert.Equal(t, ErrorTypeAuthentication, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusUnauthorized, result.Status)
}

func TestNewTokenExpiredError(t *testing.T) {
	message := "token has expired"
	underlyingErr := errors.New("jwt expired")

	result := NewTokenExpiredError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeTokenExpired, result.Code)
	assert.Equal(t, ErrorTypeAuthentication, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusUnauthorized, result.Status)
}

func TestNewAuthorizationError(t *testing.T) {
	message := "insufficient permissions"
	underlyingErr := errors.New("access denied")

	result := NewAuthorizationError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeInsufficientPermissions, result.Code)
	assert.Equal(t, ErrorTypeAuthorization, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusForbidden, result.Status)
}

func TestNewNotFoundError(t *testing.T) {
	message := "user not found"
	underlyingErr := errors.New("no rows in result set")

	result := NewNotFoundError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeResourceNotFound, result.Code)
	assert.Equal(t, ErrorTypeNotFound, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusNotFound, result.Status)
}

func TestNewConflictError(t *testing.T) {
	message := "resource already exists"
	underlyingErr := errors.New("duplicate key")

	result := NewConflictError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeResourceAlreadyExists, result.Code)
	assert.Equal(t, ErrorTypeConflict, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusConflict, result.Status)
}

func TestNewDatabaseError(t *testing.T) {
	message := "database query failed"
	underlyingErr := errors.New("connection timeout")

	result := NewDatabaseError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeDatabaseQuery, result.Code)
	assert.Equal(t, ErrorTypeDatabase, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.Status)
}

func TestNewExternalServiceError(t *testing.T) {
	message := "external service unavailable"
	underlyingErr := errors.New("service down")

	result := NewExternalServiceError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeExternalServiceDown, result.Code)
	assert.Equal(t, ErrorTypeExternal, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusBadGateway, result.Status)
}

func TestNewRateLimitError(t *testing.T) {
	message := "rate limit exceeded"
	underlyingErr := errors.New("too many requests")

	result := NewRateLimitError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeRateLimitExceeded, result.Code)
	assert.Equal(t, ErrorTypeRateLimit, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusTooManyRequests, result.Status)
}

func TestNewBusinessError(t *testing.T) {
	message := "business rule violation"
	underlyingErr := errors.New("invalid state transition")

	result := NewBusinessError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeBusinessRuleViolation, result.Code)
	assert.Equal(t, ErrorTypeBusiness, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, result.Status)
}

func TestNewInternalError(t *testing.T) {
	message := "internal server error"
	underlyingErr := errors.New("unexpected error")

	result := NewInternalError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeInternalServerError, result.Code)
	assert.Equal(t, ErrorTypeInternal, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.Status)
}

func TestNewServiceUnavailableError(t *testing.T) {
	message := "service unavailable"
	underlyingErr := errors.New("maintenance mode")

	result := NewServiceUnavailableError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeServiceUnavailable, result.Code)
	assert.Equal(t, ErrorTypeServiceUnavailable, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusServiceUnavailable, result.Status)
}

func TestNewTooManyRequestsError(t *testing.T) {
	message := "too many requests"
	underlyingErr := errors.New("rate limited")

	result := NewTooManyRequestsError(message, underlyingErr)

	assert.Equal(t, message, result.Message)
	assert.Equal(t, CodeTooManyRequests, result.Code)
	assert.Equal(t, ErrorTypeTooManyRequests, result.Type)
	assert.Equal(t, underlyingErr, result.Err)
	assert.Equal(t, http.StatusTooManyRequests, result.Status)
}

// Test legacy functions
func TestNewServiceError(t *testing.T) {
	message := "service error"
	underlyingErr := errors.New("legacy error")

	result := NewServiceError(message, underlyingErr)

	// Should behave like NewBusinessError
	expected := NewBusinessError(message, underlyingErr)
	assert.Equal(t, expected.Message, result.Message)
	assert.Equal(t, expected.Code, result.Code)
	assert.Equal(t, expected.Type, result.Type)
	assert.Equal(t, expected.Status, result.Status)
}

func TestNewBadRequestError(t *testing.T) {
	message := "bad request"
	underlyingErr := errors.New("legacy validation error")

	result := NewBadRequestError(message, underlyingErr)

	// Should behave like NewValidationError
	expected := NewValidationError(message, underlyingErr)
	assert.Equal(t, expected.Message, result.Message)
	assert.Equal(t, expected.Code, result.Code)
	assert.Equal(t, expected.Type, result.Type)
	assert.Equal(t, expected.Status, result.Status)
}

// Test error chaining
func TestAppError_Chaining(t *testing.T) {
	originalErr := errors.New("original error")
	appError := NewValidationError("validation failed", originalErr)

	// Test that the original error is preserved
	assert.Contains(t, appError.Error(), "original error")
	assert.Equal(t, originalErr, appError.Err)
}

// Test nil error handling
func TestAppError_NilError(t *testing.T) {
	tests := []struct {
		name        string
		errorFunc   func(string, error) *AppError
		expectedMsg string
	}{
		{"validation", NewValidationError, "test message"},
		{"authentication", NewAuthenticationError, "test message"},
		{"authorization", NewAuthorizationError, "test message"},
		{"not_found", NewNotFoundError, "test message"},
		{"conflict", NewConflictError, "test message"},
		{"database", NewDatabaseError, "test message"},
		{"external", NewExternalServiceError, "test message"},
		{"rate_limit", NewRateLimitError, "test message"},
		{"business", NewBusinessError, "test message"},
		{"internal", NewInternalError, "test message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errorFunc(tt.expectedMsg, nil)
			assert.Equal(t, tt.expectedMsg, result.Message)
			assert.Nil(t, result.Err)
			// Error string should not contain " - <nil>"
			assert.NotContains(t, result.Error(), " - <nil>")
		})
	}
}

// Test method chaining
func TestAppError_MethodChaining(t *testing.T) {
	originalError := NewValidationError("test error", nil)
	requestID := "req-123"
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid",
	}

	// Test chaining methods
	result := originalError.WithRequestID(requestID).WithDetails(details)

	// Should be the same instance
	assert.Same(t, originalError, result)
	assert.Equal(t, requestID, result.RequestID)
	assert.Equal(t, details, result.Details)
}

// Test constants
func TestErrorConstants(t *testing.T) {
	// Test ErrorType constants
	assert.Equal(t, ErrorType("VALIDATION_ERROR"), ErrorTypeValidation)
	assert.Equal(t, ErrorType("AUTHENTICATION_ERROR"), ErrorTypeAuthentication)
	assert.Equal(t, ErrorType("AUTHORIZATION_ERROR"), ErrorTypeAuthorization)
	assert.Equal(t, ErrorType("NOT_FOUND_ERROR"), ErrorTypeNotFound)
	assert.Equal(t, ErrorType("CONFLICT_ERROR"), ErrorTypeConflict)
	assert.Equal(t, ErrorType("DATABASE_ERROR"), ErrorTypeDatabase)
	assert.Equal(t, ErrorType("EXTERNAL_SERVICE_ERROR"), ErrorTypeExternal)
	assert.Equal(t, ErrorType("INTERNAL_ERROR"), ErrorTypeInternal)
	assert.Equal(t, ErrorType("RATE_LIMIT_ERROR"), ErrorTypeRateLimit)
	assert.Equal(t, ErrorType("BUSINESS_LOGIC_ERROR"), ErrorTypeBusiness)
	assert.Equal(t, ErrorType("SERVICE_UNAVAILABLE_ERROR"), ErrorTypeServiceUnavailable)
	assert.Equal(t, ErrorType("TOO_MANY_REQUESTS_ERROR"), ErrorTypeTooManyRequests)

	// Test error code constants
	assert.Equal(t, "E1001", CodeValidationFailed)
	assert.Equal(t, "E1002", CodeInvalidCredentials)
	assert.Equal(t, "E1003", CodeTokenExpired)
	assert.Equal(t, "E1004", CodeInsufficientPermissions)
	assert.Equal(t, "E1005", CodeResourceNotFound)
	assert.Equal(t, "E1006", CodeResourceAlreadyExists)
	assert.Equal(t, "E1007", CodeDatabaseConnection)
	assert.Equal(t, "E1008", CodeDatabaseQuery)
	assert.Equal(t, "E1009", CodeExternalServiceDown)
	assert.Equal(t, "E1010", CodeExternalServiceTimeout)
	assert.Equal(t, "E1011", CodeInternalServerError)
	assert.Equal(t, "E1012", CodeRateLimitExceeded)
	assert.Equal(t, "E1013", CodeBusinessRuleViolation)
	assert.Equal(t, "E1014", CodeServiceUnavailable)
	assert.Equal(t, "E1015", CodeTooManyRequests)
}

// Benchmark tests
func BenchmarkAppError_Error(b *testing.B) {
	appError := NewValidationError("test error", errors.New("underlying error"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = appError.Error()
	}
}

func BenchmarkNewValidationError(b *testing.B) {
	message := "validation failed"
	err := errors.New("field required")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewValidationError(message, err)
	}
}
