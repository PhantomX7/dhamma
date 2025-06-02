package utility

import (
	"context"
	"net/http"
	"testing"

	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContextWithValues(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		values ContextValues
	}{
		{
			name: "create context with complete values",
			ctx:  context.Background(),
			values: ContextValues{
				DomainID: func() *uint64 { id := uint64(123); return &id }(),
				UserID:   456,
				IsRoot:   true,
			},
		},
		{
			name: "create context with nil domain ID",
			ctx:  context.Background(),
			values: ContextValues{
				DomainID: nil,
				UserID:   789,
				IsRoot:   false,
			},
		},
		{
			name: "create context with zero values",
			ctx:  context.Background(),
			values: ContextValues{
				DomainID: func() *uint64 { id := uint64(0); return &id }(),
				UserID:   0,
				IsRoot:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultCtx := NewContextWithValues(tt.ctx, tt.values)

			// Verify the context is not nil
			require.NotNil(t, resultCtx)

			// Verify we can retrieve the values back
			retrievedValues, ok := resultCtx.Value("values").(ContextValues)
			require.True(t, ok, "should be able to retrieve ContextValues from context")

			// Verify the values match
			assert.Equal(t, tt.values.UserID, retrievedValues.UserID)
			assert.Equal(t, tt.values.IsRoot, retrievedValues.IsRoot)

			if tt.values.DomainID == nil {
				assert.Nil(t, retrievedValues.DomainID)
			} else {
				require.NotNil(t, retrievedValues.DomainID)
				assert.Equal(t, *tt.values.DomainID, *retrievedValues.DomainID)
			}
		})
	}
}

func TestValuesFromContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
		expected    ContextValues
	}{
		{
			name: "retrieve valid context values",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: func() *uint64 { id := uint64(123); return &id }(),
				UserID:   456,
				IsRoot:   true,
			}),
			expectError: false,
			expected: ContextValues{
				DomainID: func() *uint64 { id := uint64(123); return &id }(),
				UserID:   456,
				IsRoot:   true,
			},
		},
		{
			name: "retrieve context values with nil domain ID",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: nil,
				UserID:   789,
				IsRoot:   false,
			}),
			expectError: false,
			expected: ContextValues{
				DomainID: nil,
				UserID:   789,
				IsRoot:   false,
			},
		},
		{
			name:        "context without values",
			ctx:         context.Background(),
			expectError: true,
			expected:    ContextValues{},
		},
		{
			name:        "context with wrong type value",
			ctx:         context.WithValue(context.Background(), "values", "wrong_type"),
			expectError: true,
			expected:    ContextValues{},
		},
		{
			name:        "context with nil value",
			ctx:         context.WithValue(context.Background(), "values", nil),
			expectError: true,
			expected:    ContextValues{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := ValuesFromContext(tt.ctx)

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, &errors.AppError{}, err)

				appErr := err.(*errors.AppError)
				assert.Equal(t, "forbidden", appErr.Message)
				assert.Equal(t, http.StatusBadRequest, appErr.Status)

				// Should return empty ContextValues on error
				assert.Equal(t, ContextValues{}, values)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.expected.UserID, values.UserID)
				assert.Equal(t, tt.expected.IsRoot, values.IsRoot)

				if tt.expected.DomainID == nil {
					assert.Nil(t, values.DomainID)
				} else {
					require.NotNil(t, values.DomainID)
					assert.Equal(t, *tt.expected.DomainID, *values.DomainID)
				}
			}
		})
	}
}

func TestCheckDomainContext(t *testing.T) {
	tests := []struct {
		name            string
		ctx             context.Context
		entityDomainID  uint64
		entityName      string
		actionVerb      string
		expectError     bool
		expectedMessage string
		expectedValues  ContextValues
	}{
		{
			name: "valid domain match",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: func() *uint64 { id := uint64(123); return &id }(),
				UserID:   456,
				IsRoot:   true,
			}),
			entityDomainID: 123,
			entityName:     "product",
			actionVerb:     "create",
			expectError:    false,
			expectedValues: ContextValues{
				DomainID: func() *uint64 { id := uint64(123); return &id }(),
				UserID:   456,
				IsRoot:   true,
			},
		},
		{
			name: "nil domain ID in context - should allow operation",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: nil,
				UserID:   789,
				IsRoot:   false,
			}),
			entityDomainID: 456,
			entityName:     "follower",
			actionVerb:     "update",
			expectError:    false,
			expectedValues: ContextValues{
				DomainID: nil,
				UserID:   789,
				IsRoot:   false,
			},
		},
		{
			name: "domain mismatch",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: func() *uint64 { id := uint64(123); return &id }(),
				UserID:   456,
				IsRoot:   false,
			}),
			entityDomainID:  456,
			entityName:      "product",
			actionVerb:      "delete",
			expectError:     true,
			expectedMessage: "you cannot delete product for another domain",
		},
		{
			name:            "context without values",
			ctx:             context.Background(),
			entityDomainID:  123,
			entityName:      "user",
			actionVerb:      "get",
			expectError:     true,
			expectedMessage: "forbidden",
		},
		{
			name: "zero domain IDs match",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: func() *uint64 { id := uint64(0); return &id }(),
				UserID:   100,
				IsRoot:   false,
			}),
			entityDomainID: 0,
			entityName:     "category",
			actionVerb:     "create",
			expectError:    false,
			expectedValues: ContextValues{
				DomainID: func() *uint64 { id := uint64(0); return &id }(),
				UserID:   100,
				IsRoot:   false,
			},
		},
		{
			name: "different entity names and verbs",
			ctx: NewContextWithValues(context.Background(), ContextValues{
				DomainID: func() *uint64 { id := uint64(100); return &id }(),
				UserID:   200,
				IsRoot:   true,
			}),
			entityDomainID:  200,
			entityName:      "order",
			actionVerb:      "process",
			expectError:     true,
			expectedMessage: "you cannot process order for another domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := CheckDomainContext(tt.ctx, tt.entityDomainID, tt.entityName, tt.actionVerb)

			if tt.expectError {
				require.Error(t, err)
				assert.IsType(t, &errors.AppError{}, err)

				appErr := err.(*errors.AppError)
				assert.Equal(t, tt.expectedMessage, appErr.Message)
				assert.Equal(t, http.StatusBadRequest, appErr.Status)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.expectedValues.UserID, values.UserID)
				assert.Equal(t, tt.expectedValues.IsRoot, values.IsRoot)

				if tt.expectedValues.DomainID == nil {
					assert.Nil(t, values.DomainID)
				} else {
					require.NotNil(t, values.DomainID)
					assert.Equal(t, *tt.expectedValues.DomainID, *values.DomainID)
				}
			}
		})
	}
}

// TestContextValues_EdgeCases tests edge cases for ContextValues struct
func TestContextValues_EdgeCases(t *testing.T) {
	t.Run("empty context values", func(t *testing.T) {
		emptyValues := ContextValues{}
		ctx := NewContextWithValues(context.Background(), emptyValues)

		retrieved, err := ValuesFromContext(ctx)
		require.NoError(t, err)

		assert.Nil(t, retrieved.DomainID)
		assert.Equal(t, uint64(0), retrieved.UserID)
		assert.False(t, retrieved.IsRoot)
	})

	t.Run("large domain and user IDs", func(t *testing.T) {
		largeID := uint64(18446744073709551615) // max uint64
		values := ContextValues{
			DomainID: &largeID,
			UserID:   largeID,
			IsRoot:   true,
		}

		ctx := NewContextWithValues(context.Background(), values)
		retrieved, err := ValuesFromContext(ctx)

		require.NoError(t, err)
		require.NotNil(t, retrieved.DomainID)
		assert.Equal(t, largeID, *retrieved.DomainID)
		assert.Equal(t, largeID, retrieved.UserID)
		assert.True(t, retrieved.IsRoot)
	})
}

// BenchmarkNewContextWithValues benchmarks the context creation
func BenchmarkNewContextWithValues(b *testing.B) {
	values := ContextValues{
		DomainID: func() *uint64 { id := uint64(123); return &id }(),
		UserID:   456,
		IsRoot:   true,
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewContextWithValues(ctx, values)
	}
}

// BenchmarkValuesFromContext benchmarks the context value retrieval
func BenchmarkValuesFromContext(b *testing.B) {
	values := ContextValues{
		DomainID: func() *uint64 { id := uint64(123); return &id }(),
		UserID:   456,
		IsRoot:   true,
	}
	ctx := NewContextWithValues(context.Background(), values)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValuesFromContext(ctx)
	}
}
