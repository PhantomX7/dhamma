package pagination

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestScopeBuilder_StressTests tests performance and memory usage under stress
func TestScopeBuilder_StressTests(t *testing.T) {
	t.Run("large number of filter conditions", func(t *testing.T) {
		filterDef := NewFilterDefinition()
		conditions := make(map[string][]string)

		// Create 100 different filters
		for i := 0; i < 100; i++ {
			fieldName := fmt.Sprintf("field_%d", i)
			filterDef.AddFilter(fieldName, FilterConfig{
				Field:     fieldName,
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorLike},
			})
			conditions[fieldName] = []string{fmt.Sprintf("value_%d", i)}
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)

		start := time.Now()
		filterScopes, metaScopes := scopeBuilder.Build()
		duration := time.Since(start)

		// Should complete within reasonable time (less than 100ms)
		assert.Less(t, duration, 100*time.Millisecond)
		assert.Equal(t, 100, len(filterScopes))
		assert.Equal(t, 2, len(metaScopes)) // limit and order
	})

	t.Run("large number of values in IN operation", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("id", FilterConfig{
			Field:     "id",
			Type:      FilterTypeNumber,
			Operators: []FilterOperator{OperatorIn},
		})

		// Create IN condition with 1000 values
		values := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			values[i] = fmt.Sprintf("%d", i)
		}
		inValue := "in:" + strings.Join(values, ",")

		conditions := map[string][]string{
			"id": {inValue},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)

		start := time.Now()
		filterScopes, _ := scopeBuilder.Build()
		duration := time.Since(start)

		// Should complete within reasonable time
		assert.Less(t, duration, 50*time.Millisecond)
		assert.Equal(t, 1, len(filterScopes))
	})

	t.Run("deeply nested search fields", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("search", FilterConfig{
			Field:        "search",
			SearchFields: make([]string, 50), // 50 search fields
			Type:         FilterTypeString,
			Operators:    []FilterOperator{OperatorLike},
		})

		// Create search fields slice
		searchFields := make([]string, 50)
		for i := 0; i < 50; i++ {
			searchFields[i] = fmt.Sprintf("field_%d", i)
		}
		
		// Update the filter config with search fields
		filterDef.AddFilter("search", FilterConfig{
			Field:        "search",
			SearchFields: searchFields,
			Type:         FilterTypeString,
			Operators:    []FilterOperator{OperatorLike},
		})

		conditions := map[string][]string{
			"search": {"like:test"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)

		start := time.Now()
		filterScopes, _ := scopeBuilder.Build()
		duration := time.Since(start)

		// Should complete within reasonable time
		assert.Less(t, duration, 50*time.Millisecond)
		assert.Equal(t, 1, len(filterScopes))
	})
}

// TestScopeBuilder_ConcurrencyTests tests thread safety
func TestScopeBuilder_ConcurrencyTests(t *testing.T) {
	t.Run("concurrent scope building", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorEquals},
		})

		conditions := map[string][]string{
			"name": {"test"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})

		// Run 100 goroutines concurrently
		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func(id int) {
				scopeBuilder := NewScopeBuilder(pagination)
				filterScopes, metaScopes := scopeBuilder.Build()
				
				// Basic assertions
				assert.Equal(t, 1, len(filterScopes))
				assert.Equal(t, 2, len(metaScopes))
				
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 100; i++ {
			select {
			case <-done:
				// Success
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout waiting for concurrent operations")
			}
		}
	})
}

// TestScopeBuilder_MemoryLeakTests tests for potential memory leaks
func TestScopeBuilder_MemoryLeakTests(t *testing.T) {
	t.Run("repeated scope building", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorEquals},
		})

		conditions := map[string][]string{
			"name": {"test"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})

		// Build scopes many times to check for memory leaks
		for i := 0; i < 1000; i++ {
			scopeBuilder := NewScopeBuilder(pagination)
			filterScopes, metaScopes := scopeBuilder.Build()
			
			// Verify scopes are created correctly
			assert.Equal(t, 1, len(filterScopes))
			assert.Equal(t, 2, len(metaScopes))
			
			// Clear references to help GC
			filterScopes = nil
			metaScopes = nil
			scopeBuilder = nil
		}
	})
}

// TestScopeBuilder_ComplexScenarios tests complex real-world scenarios
func TestScopeBuilder_ComplexScenarios(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the test model
	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	dryRunDB := db.Session(&gorm.Session{DryRun: true})

	t.Run("e-commerce product search scenario", func(t *testing.T) {
		// Simulate complex e-commerce product filtering
		filterDef := NewFilterDefinition().
			AddFilter("search", FilterConfig{
				Field:        "search",
				SearchFields: []string{"name", "description", "brand", "category"},
				Type:         FilterTypeString,
				Operators:    []FilterOperator{OperatorLike},
			}).
			AddFilter("price", FilterConfig{
				Field:     "price",
				Type:      FilterTypeNumber,
				Operators: []FilterOperator{OperatorBetween, OperatorGte, OperatorLte},
			}).
			AddFilter("category", FilterConfig{
				Field:      "category_id",
				Type:       FilterTypeEnum,
				Operators:  []FilterOperator{OperatorIn, OperatorEquals},
				EnumValues: []string{"electronics", "clothing", "books", "home"},
			}).
			AddFilter("available", FilterConfig{
				Field:     "is_available",
				Type:      FilterTypeBool,
				Operators: []FilterOperator{OperatorEquals},
			}).
			AddFilter("created_date", FilterConfig{
				Field:     "created_at",
				Type:      FilterTypeDate,
				Operators: []FilterOperator{OperatorBetween, OperatorGte},
			})

		conditions := map[string][]string{
			"search":       {"like:laptop"},
			"price":        {"between:100,1000"},
			"category":     {"in:electronics,home"},
			"available":    {"eq:true"},
			"created_date": {"gte:2023-01-01"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{
			DefaultLimit: 20,
			MaxLimit:     100,
			DefaultOrder: "created_at DESC",
		})

		scopeBuilder := NewScopeBuilder(pagination)
		filterScopes, metaScopes := scopeBuilder.Build()

		// Should create all filter scopes
		assert.Equal(t, 5, len(filterScopes))
		assert.Equal(t, 2, len(metaScopes))

		// Test SQL generation
		var models []TestModel
		query := dryRunDB.Model(&TestModel{})

		// Apply all scopes
		for _, scope := range filterScopes {
			query = scope(query)
		}
		for _, scope := range metaScopes {
			query = scope(query)
		}

		stmt := query.Find(&models).Statement
		sql := strings.ToUpper(stmt.SQL.String())

		// Verify complex SQL is generated correctly
		assert.Contains(t, sql, "LIKE")
		assert.Contains(t, sql, "BETWEEN")
		assert.Contains(t, sql, "IN")
		assert.Contains(t, sql, "LIMIT")
		assert.Contains(t, sql, "ORDER BY")
	})

	t.Run("user management filtering scenario", func(t *testing.T) {
		// Simulate user management system filtering
		filterDef := NewFilterDefinition().
			AddFilter("name", FilterConfig{
				Field:        "name",
				SearchFields: []string{"first_name", "last_name", "username", "email"},
				Type:         FilterTypeString,
				Operators:    []FilterOperator{OperatorLike, OperatorEquals},
			}).
			AddFilter("role", FilterConfig{
				Field:      "role",
				Type:       FilterTypeEnum,
				Operators:  []FilterOperator{OperatorIn, OperatorEquals},
				EnumValues: []string{"admin", "user", "moderator", "guest"},
			}).
			AddFilter("active", FilterConfig{
				Field:     "is_active",
				Type:      FilterTypeBool,
				Operators: []FilterOperator{OperatorEquals},
			}).
			AddFilter("last_login", FilterConfig{
				Field:     "last_login_at",
				Type:      FilterTypeDate,
				Operators: []FilterOperator{OperatorGte, OperatorLte, OperatorBetween},
			}).
			AddFilter("age", FilterConfig{
				Field:     "age",
				Type:      FilterTypeNumber,
				Operators: []FilterOperator{OperatorGte, OperatorLte, OperatorBetween},
			})

		conditions := map[string][]string{
			"name":       {"like:john"},
			"role":       {"in:admin,moderator"},
			"active":     {"eq:true"},
			"last_login": {"gte:2023-01-01"},
			"age":        {"between:18,65"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{
			DefaultLimit: 50,
			MaxLimit:     200,
			DefaultOrder: "last_login_at DESC, created_at DESC",
		})

		scopeBuilder := NewScopeBuilder(pagination)
		filterScopes, metaScopes := scopeBuilder.Build()

		// Should create all filter scopes
		assert.Equal(t, 5, len(filterScopes))
		assert.Equal(t, 2, len(metaScopes))

		// Test that scopes can be applied without errors
		var models []TestModel
		query := dryRunDB.Model(&TestModel{})

		for _, scope := range filterScopes {
			query = scope(query)
		}
		for _, scope := range metaScopes {
			query = scope(query)
		}

		stmt := query.Find(&models).Statement
		assert.NotEmpty(t, stmt.SQL.String())
		assert.NotEmpty(t, stmt.Vars)
	})
}

// TestScopeBuilder_ErrorRecovery tests error recovery scenarios
func TestScopeBuilder_ErrorRecovery(t *testing.T) {
	t.Run("malformed filter values", func(t *testing.T) {
		filterDef := NewFilterDefinition().
			AddFilter("number", FilterConfig{
				Field:     "number_field",
				Type:      FilterTypeNumber,
				Operators: []FilterOperator{OperatorBetween},
			}).
			AddFilter("date", FilterConfig{
				Field:     "date_field",
				Type:      FilterTypeDate,
				Operators: []FilterOperator{OperatorEquals},
			})

		// Malformed conditions that should be gracefully handled
		conditions := map[string][]string{
			"number": {"invalid:123"}, // using disallowed operator for number
			"date":   {"like:invalid-date-format"}, // using disallowed operator for date
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)

		// Should not panic and should handle errors gracefully
		assert.NotPanics(t, func() {
			filterScopes, metaScopes := scopeBuilder.Build()
			
			// Should not create scopes for malformed data
			assert.Equal(t, 0, len(filterScopes))
			assert.Equal(t, 2, len(metaScopes)) // meta scopes should still be created
		})
	})

	t.Run("extremely long filter values", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike},
		})

		// Create extremely long string (10KB)
		longString := strings.Repeat("a", 10*1024)
		conditions := map[string][]string{
			"name": {"like:" + longString},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)

		// Should handle long strings without issues
		assert.NotPanics(t, func() {
			filterScopes, _ := scopeBuilder.Build()
			assert.Equal(t, 1, len(filterScopes))
		})
	})

	t.Run("unicode and special characters", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorEquals, OperatorLike},
		})

		// Test with various unicode and special characters
		specialValues := []string{
			"like:测试",                    // Chinese characters
			"eq:🚀🎉",                     // Emojis
			"like:Ñoño",                  // Accented characters
			"eq:user@domain.com",         // Email format
			"like:O'Connor",              // Apostrophe
			"eq:100%",                    // Percentage
			"like:$1000.00",              // Currency format
			"eq:file.txt",                // Dot notation
			"like:path/to/file",          // Path separators
			"eq:key=value&other=value2",  // Query string format
		}

		for i, value := range specialValues {
			t.Run(fmt.Sprintf("special_char_%d", i), func(t *testing.T) {
				conditions := map[string][]string{
					"name": {value},
				}

				pagination := NewPagination(conditions, filterDef, PaginationOptions{})
				scopeBuilder := NewScopeBuilder(pagination)

				assert.NotPanics(t, func() {
					filterScopes, _ := scopeBuilder.Build()
					assert.Equal(t, 1, len(filterScopes))
				})
			})
		}
	})
}

// BenchmarkScopeBuilder_EdgeCases benchmarks edge case performance
func BenchmarkScopeBuilder_EdgeCases(b *testing.B) {
	b.Run("large_filter_set", func(b *testing.B) {
		filterDef := NewFilterDefinition()
		conditions := make(map[string][]string)

		// Create 50 filters
		for i := 0; i < 50; i++ {
			fieldName := fmt.Sprintf("field_%d", i)
			filterDef.AddFilter(fieldName, FilterConfig{
				Field:     fieldName,
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals},
			})
			conditions[fieldName] = []string{fmt.Sprintf("value_%d", i)}
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			scopeBuilder := NewScopeBuilder(pagination)
			filterScopes, metaScopes := scopeBuilder.Build()
			_ = filterScopes
			_ = metaScopes
		}
	})

	b.Run("complex_string_search", func(b *testing.B) {
		filterDef := NewFilterDefinition().AddFilter("search", FilterConfig{
			Field:        "search",
			SearchFields: []string{"name", "description", "tags", "category", "brand", "model", "sku", "title", "content", "summary"},
			Type:         FilterTypeString,
			Operators:    []FilterOperator{OperatorLike},
		})

		conditions := map[string][]string{
			"search": {"like:complex search term with multiple words"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			scopeBuilder := NewScopeBuilder(pagination)
			filterScopes, metaScopes := scopeBuilder.Build()
			_ = filterScopes
			_ = metaScopes
		}
	})

	b.Run("large_in_operation", func(b *testing.B) {
		filterDef := NewFilterDefinition().AddFilter("ids", FilterConfig{
			Field:     "id",
			Type:      FilterTypeNumber,
			Operators: []FilterOperator{OperatorIn},
		})

		// Create IN with 500 values
		values := make([]string, 500)
		for i := 0; i < 500; i++ {
			values[i] = fmt.Sprintf("%d", i)
		}
		inValue := "in:" + strings.Join(values, ",")

		conditions := map[string][]string{
			"ids": {inValue},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			scopeBuilder := NewScopeBuilder(pagination)
			filterScopes, metaScopes := scopeBuilder.Build()
			_ = filterScopes
			_ = metaScopes
		}
	})
}