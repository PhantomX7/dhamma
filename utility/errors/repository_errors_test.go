package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		format   string
		args     []any
		expected string
	}{
		{
			name:     "wrap with format",
			err:      errors.New("original error"),
			format:   "failed to process %s",
			args:     []any{"user data"},
			expected: "failed to process user data",
		},
		{
			name:     "wrap without args",
			err:      errors.New("database error"),
			format:   "database operation failed",
			args:     nil,
			expected: "database operation failed",
		},
		{
			name:     "nil error",
			err:      nil,
			format:   "should not wrap",
			args:     nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.format, tt.args...)
			if tt.err == nil {
				assert.Nil(t, result)
			} else {
				appErr, ok := result.(*AppError)
				assert.True(t, ok)
				assert.Equal(t, tt.expected, appErr.Message)
				assert.Equal(t, tt.err, appErr.Err)
				assert.Equal(t, http.StatusUnprocessableEntity, appErr.Status)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			expected: true,
		},
		{
			name:     "other error",
			err:      ErrDuplicate,
			expected: false,
		},
		{
			name:     "custom error",
			err:      errors.New("custom error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrDuplicate",
			err:      ErrDuplicate,
			expected: true,
		},
		{
			name:     "other error",
			err:      ErrNotFound,
			expected: false,
		},
		{
			name:     "custom error",
			err:      errors.New("custom error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDuplicate(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDatabaseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrDatabase",
			err:      ErrDatabase,
			expected: true,
		},
		{
			name:     "other error",
			err:      ErrNotFound,
			expected: false,
		},
		{
			name:     "custom error",
			err:      errors.New("custom error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDatabaseError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test error constants
func TestRepositoryErrorConstants(t *testing.T) {
	// Test that all error constants are properly defined
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrDuplicate)
	assert.NotNil(t, ErrInvalidData)
	assert.NotNil(t, ErrForeignKey)
	assert.NotNil(t, ErrDatabase)
	assert.NotNil(t, ErrPermissionDenied)
	assert.NotNil(t, ErrInvalidRequest)

	// Test error messages
	assert.Equal(t, "record not found", ErrNotFound.Error())
	assert.Equal(t, "duplicate record", ErrDuplicate.Error())
	assert.Equal(t, "invalid data", ErrInvalidData.Error())
	assert.Equal(t, "foreign key constraint failed", ErrForeignKey.Error())
	assert.Equal(t, "database error", ErrDatabase.Error())
	assert.Equal(t, "permission denied", ErrPermissionDenied.Error())
	assert.Equal(t, "invalid request", ErrInvalidRequest.Error())
}

// Test error wrapping behavior
func TestWrapError_Preserves_Original(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, "context: %s", "test")

	appErr, ok := wrappedErr.(*AppError)
	assert.True(t, ok)
	assert.Equal(t, "context: test", appErr.Message)
	assert.Equal(t, originalErr, appErr.Err)
	assert.Equal(t, http.StatusUnprocessableEntity, appErr.Status)
}

// Test multiple wrapping behavior
func TestWrapError_Multiple_Wrapping(t *testing.T) {
	originalErr := ErrNotFound
	firstWrap := WrapError(originalErr, "first wrap")
	secondWrap := WrapError(firstWrap, "second wrap")

	// Check the structure
	firstAppErr, ok := firstWrap.(*AppError)
	assert.True(t, ok)
	assert.Equal(t, "first wrap", firstAppErr.Message)
	assert.Equal(t, originalErr, firstAppErr.Err)

	secondAppErr, ok := secondWrap.(*AppError)
	assert.True(t, ok)
	assert.Equal(t, "second wrap", secondAppErr.Message)
	assert.Equal(t, firstWrap, secondAppErr.Err)

	// Test that we can access the underlying error
	assert.Equal(t, originalErr, firstAppErr.Err)
	assert.Equal(t, firstWrap, secondAppErr.Err)
}

// Test error chaining with errors.Is
func TestErrorChaining(t *testing.T) {
	tests := []struct {
		name      string
		baseError error
		checkFunc func(error) bool
		expected  bool
	}{
		{
			name:      "direct not found",
			baseError: ErrNotFound,
			checkFunc: IsNotFound,
			expected:  true,
		},
		{
			name:      "direct duplicate",
			baseError: ErrDuplicate,
			checkFunc: IsDuplicate,
			expected:  true,
		},
		{
			name:      "direct database",
			baseError: ErrDatabase,
			checkFunc: IsDatabaseError,
			expected:  true,
		},
		{
			name:      "wrong error type",
			baseError: ErrNotFound,
			checkFunc: IsDuplicate,
			expected:  false,
		},
		{
			name:      "custom error",
			baseError: errors.New("custom error"),
			checkFunc: IsNotFound,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checkFunc(tt.baseError)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkWrapError(b *testing.B) {
	err := errors.New("test error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WrapError(err, "context: %s", "test")
	}
}

func BenchmarkIsNotFound(b *testing.B) {
	err := ErrNotFound
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsNotFound(err)
	}
}

func BenchmarkIsDuplicate(b *testing.B) {
	err := ErrDuplicate
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsDuplicate(err)
	}
}

func BenchmarkIsDatabaseError(b *testing.B) {
	err := ErrDatabase
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsDatabaseError(err)
	}
}
