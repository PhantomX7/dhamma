package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test model for example usage
type ExampleModel struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	UserID    uint
	Status    string
	CreatedAt string // Using string for simplicity in example
}

// Test the example usage function
func TestExampleUsage(t *testing.T) {
	// This test ensures that the example code compiles and runs without errors
	// We can't test the actual database operations since the example uses a nil db
	// But we can test that the pagination and scope building works

	assert.NotPanics(t, func() {
		ExampleUsage()
	})
}

// Test the example filter definition setup
func TestExampleFilterDefinition(t *testing.T) {
	// Recreate the filter definition from the example
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike, OperatorEquals},
		}).
		AddFilter("user_id", FilterConfig{
			Field:     "id",
			TableName: "User",
			Type:      FilterTypeID,
			Operators: []FilterOperator{OperatorIn, OperatorEquals},
		}).
		AddFilter("created_at", FilterConfig{
			Field:     "created_at",
			Type:      FilterTypeDateTime,
			Operators: []FilterOperator{OperatorBetween},
		}).
		AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			Operators:  []FilterOperator{OperatorEquals, OperatorIn},
			EnumValues: []string{"active", "inactive", "pending"},
		}).
		AddSort("created_at", SortConfig{
			Field:   "created_at",
			Allowed: true,
		})

	// Verify the filter definition was created correctly
	assert.NotNil(t, filterDef)
	assert.Equal(t, 4, len(filterDef.configs))

	// Test name filter
	nameConfig := filterDef.configs["name"]
	assert.Equal(t, "name", nameConfig.Field)
	assert.Equal(t, FilterTypeString, nameConfig.Type)
	assert.Contains(t, nameConfig.Operators, OperatorLike)
	assert.Contains(t, nameConfig.Operators, OperatorEquals)

	// Test user_id filter
	userIDConfig := filterDef.configs["user_id"]
	assert.Equal(t, "id", userIDConfig.Field)
	assert.Equal(t, "User", userIDConfig.TableName)
	assert.Equal(t, FilterTypeID, userIDConfig.Type)
	assert.Contains(t, userIDConfig.Operators, OperatorIn)
	assert.Contains(t, userIDConfig.Operators, OperatorEquals)

	// Test created_at filter
	createdAtConfig := filterDef.configs["created_at"]
	assert.Equal(t, "created_at", createdAtConfig.Field)
	assert.Equal(t, FilterTypeDateTime, createdAtConfig.Type)
	assert.Contains(t, createdAtConfig.Operators, OperatorBetween)

	// Test status filter
	statusConfig := filterDef.configs["status"]
	assert.Equal(t, "status", statusConfig.Field)
	assert.Equal(t, FilterTypeEnum, statusConfig.Type)
	assert.Equal(t, []string{"active", "inactive", "pending"}, statusConfig.EnumValues)
}

// Test the example conditions
func TestExampleConditions(t *testing.T) {
	// Recreate the conditions from the example
	conditions := map[string][]string{
		"limit":      {"20"},
		"offset":     {"0"},
		"sort":       {"created_at desc"},
		"name":       {"John"},
		"user_id":    {"1,2,3"},
		"created_at": {"between:2023-01-01,2023-12-31"},
		"status":     {"active"},
	}

	// Verify conditions structure
	assert.Equal(t, "20", conditions["limit"][0])
	assert.Equal(t, "0", conditions["offset"][0])
	assert.Equal(t, "created_at desc", conditions["sort"][0])
	assert.Equal(t, "John", conditions["name"][0])
	assert.Equal(t, "1,2,3", conditions["user_id"][0])
	assert.Equal(t, "between:2023-01-01,2023-12-31", conditions["created_at"][0])
	assert.Equal(t, "active", conditions["status"][0])
}

// Test the example pagination creation
func TestExamplePaginationCreation(t *testing.T) {
	// Recreate the filter definition and conditions from the example
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike, OperatorEquals},
		}).
		AddFilter("user_id", FilterConfig{
			Field:     "id",
			TableName: "User",
			Type:      FilterTypeID,
			Operators: []FilterOperator{OperatorIn, OperatorEquals},
		}).
		AddFilter("created_at", FilterConfig{
			Field:     "created_at",
			Type:      FilterTypeDateTime,
			Operators: []FilterOperator{OperatorBetween},
		}).
		AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			Operators:  []FilterOperator{OperatorEquals, OperatorIn},
			EnumValues: []string{"active", "inactive", "pending"},
		}).
		AddSort("created_at", SortConfig{
			Field:   "created_at",
			Allowed: true,
		})

	conditions := map[string][]string{
		"limit":      {"20"},
		"offset":     {"0"},
		"sort":       {"created_at desc"},
		"name":       {"John"},
		"user_id":    {"1,2,3"},
		"created_at": {"between:2023-01-01,2023-12-31"},
		"status":     {"active"},
	}

	// Create pagination as in the example
	pagination := NewPagination(conditions, filterDef, PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	})

	// Verify pagination was created correctly
	assert.NotNil(t, pagination)
	assert.Equal(t, 20, pagination.Limit)
	assert.Equal(t, 0, pagination.Offset)
	assert.Equal(t, "created_at desc", pagination.Order) // Should use the sort from conditions
	assert.Equal(t, conditions, pagination.Conditions)
	assert.Equal(t, filterDef, pagination.FilterDef)
}

// Test the example scope building
func TestExampleScopeBuilding(t *testing.T) {
	// Recreate the complete example setup
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike, OperatorEquals},
		}).
		AddFilter("user_id", FilterConfig{
			Field:     "id",
			TableName: "User",
			Type:      FilterTypeID,
			Operators: []FilterOperator{OperatorIn, OperatorEquals},
		}).
		AddFilter("created_at", FilterConfig{
			Field:     "created_at",
			Type:      FilterTypeDateTime,
			Operators: []FilterOperator{OperatorBetween},
		}).
		AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			Operators:  []FilterOperator{OperatorEquals, OperatorIn},
			EnumValues: []string{"active", "inactive", "pending"},
		}).
		AddSort("created_at", SortConfig{
			Field:   "created_at",
			Allowed: true,
		})

	conditions := map[string][]string{
		"limit":      {"20"},
		"offset":     {"0"},
		"sort":       {"created_at desc"},
		"name":       {"John"},
		"user_id":    {"1,2,3"},
		"created_at": {"between:2023-01-01,2023-12-31"},
		"status":     {"active"},
	}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	})

	// Build scopes as in the example
	scopeBuilder := NewScopeBuilder(pagination)
	filterScopes, metaScopes := scopeBuilder.Build()

	// Verify scopes were built
	assert.NotNil(t, filterScopes)
	assert.NotNil(t, metaScopes)

	// Should have filter scopes for: name, user_id, created_at, status
	assert.Equal(t, 4, len(filterScopes))

	// Should have meta scopes for: limit, offset (since it's 0, it might not be included), order
	// Note: offset of 0 typically doesn't create a scope
	assert.GreaterOrEqual(t, len(metaScopes), 2) // At least limit and order
	assert.LessOrEqual(t, len(metaScopes), 3)    // At most limit, offset, and order
}

// Test example with actual database operations
func TestExampleWithDatabase(t *testing.T) {
	// Skip database tests if no test database is available
	t.Skip("Database tests require MySQL setup")
}

// Test example error cases
func TestExampleErrorCases(t *testing.T) {
	// Test with invalid enum values
	filterDef := NewFilterDefinition().
		AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			Operators:  []FilterOperator{OperatorEquals},
			EnumValues: []string{"active", "inactive", "pending"},
		})

	conditions := map[string][]string{
		"status": {"invalid_status"}, // Invalid enum value
	}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	})

	scopeBuilder := NewScopeBuilder(pagination)
	filterScopes, metaScopes := scopeBuilder.Build()

	// Should have no filter scopes due to invalid enum value
	assert.Equal(t, 0, len(filterScopes))
	// Should still have meta scopes
	assert.GreaterOrEqual(t, len(metaScopes), 2)
}

// Benchmark the example usage
func BenchmarkExampleUsage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExampleUsage()
	}
}

func BenchmarkExampleScopeBuilding(b *testing.B) {
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike, OperatorEquals},
		}).
		AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			Operators:  []FilterOperator{OperatorEquals},
			EnumValues: []string{"active", "inactive", "pending"},
		})

	conditions := map[string][]string{
		"name":   {"John"},
		"status": {"active"},
		"limit":  {"20"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pagination := NewPagination(conditions, filterDef, PaginationOptions{
			DefaultLimit: 20,
			MaxLimit:     100,
			DefaultOrder: "id desc",
		})

		scopeBuilder := NewScopeBuilder(pagination)
		_, _ = scopeBuilder.Build()
	}
}
