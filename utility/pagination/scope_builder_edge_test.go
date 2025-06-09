package pagination

import (
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestScopeBuilder_EdgeCases tests various edge cases and error conditions
func TestScopeBuilder_EdgeCases(t *testing.T) {
	t.Run("nil pagination", func(t *testing.T) {
		// Test with nil pagination - should not panic
		sb := NewScopeBuilder(nil)
		assert.NotNil(t, sb)
		assert.Nil(t, sb.pagination)
		assert.NotNil(t, sb.scopes)
		assert.NotNil(t, sb.metaScopes)
	})

	t.Run("empty conditions", func(t *testing.T) {
		filterDef := NewFilterDefinition()
		conditions := map[string][]string{}
		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		sb := NewScopeBuilder(pagination)
		filterScopes, metaScopes := sb.Build()
		
		assert.NotNil(t, filterScopes)
		assert.NotNil(t, metaScopes)
		assert.Equal(t, 0, len(filterScopes))
		assert.Equal(t, 2, len(metaScopes)) // limit and order
	})

	t.Run("unknown filter field", func(t *testing.T) {
		filterDef := NewFilterDefinition()
		conditions := map[string][]string{
			"unknown_field": {"test"},
		}
		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		sb := NewScopeBuilder(pagination)
		filterScopes, _ := sb.Build()
		
		// Should ignore unknown fields
		assert.Equal(t, 0, len(filterScopes))
	})

	t.Run("empty values array", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorEquals},
		})
		conditions := map[string][]string{
			"name": {}, // empty values
		}
		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		sb := NewScopeBuilder(pagination)
		filterScopes, _ := sb.Build()
		
		// Should not create scope for empty values
		assert.Equal(t, 0, len(filterScopes))
	})

	t.Run("invalid operator for filter type", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("age", FilterConfig{
			Field:     "age",
			Type:      FilterTypeNumber,
			Operators: []FilterOperator{OperatorEquals}, // only equals allowed
		})
		conditions := map[string][]string{
			"age": {"like:25"}, // like not allowed for numbers
		}
		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		sb := NewScopeBuilder(pagination)
		filterScopes, _ := sb.Build()
		
		// Should not create scope for invalid operator
		assert.Equal(t, 0, len(filterScopes))
	})
}

// TestParseFilterOperation_EdgeCases tests edge cases in filter operation parsing
func TestParseFilterOperation_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected FilterOperation
	}{
		{
			name:  "empty string",
			value: "",
			expected: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{""},
			},
		},
		{
			name:  "only colon",
			value: ":",
			expected: FilterOperation{
				Operator: FilterOperator(""),
				Values:   []string{""},
			},
		},
		{
			name:  "operator with empty value",
			value: "eq:",
			expected: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{""},
			},
		},
		{
			name:  "in operator with empty values",
			value: "in:",
			expected: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{""},
			},
		},
		{
			name:  "in operator with comma only",
			value: "in:,",
			expected: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"", ""},
			},
		},
		{
			name:  "multiple consecutive colons",
			value: "eq::value",
			expected: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{":value"},
			},
		},
		{
			name:  "special characters in value",
			value: "eq:test@#$%^&*()",
			expected: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"test@#$%^&*()"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFilterOperation(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsValidOperation_EdgeCases tests edge cases in operation validation
func TestIsValidOperation_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		operation FilterOperation
		config    FilterConfig
		expected  bool
	}{
		{
			name: "empty operators list",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"test"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{}, // empty
			},
			expected: false,
		},
		{
			name: "between with three values",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"1", "2", "3"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorBetween},
			},
			expected: false,
		},
		{
			name: "between with empty values",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorBetween},
			},
			expected: false,
		},
		{
			name: "not_in with empty values",
			operation: FilterOperation{
				Operator: OperatorNotIn,
				Values:   []string{},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorNotIn},
			},
			expected: false,
		},
		{
			name: "equals with zero values",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorEquals},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidOperation(tt.operation, tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBuildNumberScope_EdgeCases tests edge cases in number scope building
func TestBuildNumberScope_EdgeCases(t *testing.T) {
	sb := &ScopeBuilder{}

	tests := []struct {
		name      string
		field     string
		operation FilterOperation
		expected  bool // whether scope should be created
	}{
		{
			name:  "empty field name",
			field: "",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"25"},
			},
			expected: true, // should still create scope
		},
		{
			name:  "unsupported operator",
			field: "age",
			operation: FilterOperation{
				Operator: OperatorLike, // not supported for numbers
				Values:   []string{"25"},
			},
			expected: false,
		},
		{
			name:  "empty operator",
			field: "age",
			operation: FilterOperation{
				Operator: FilterOperator(""),
				Values:   []string{"25"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildNumberScope(tt.field, tt.operation)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// TestBuildStringScope_EdgeCases tests edge cases in string scope building
func TestBuildStringScope_EdgeCases(t *testing.T) {
	sb := &ScopeBuilder{}

	tests := []struct {
		name      string
		fields    []string
		operation FilterOperation
		expected  bool // whether scope should be created
	}{
		{
			name:   "empty fields array",
			fields: []string{},
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"test"},
			},
			expected: false,
		},
		{
			name:   "nil fields array",
			fields: nil,
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"test"},
			},
			expected: false,
		},
		{
			name:   "empty field name in array",
			fields: []string{""},
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"test"},
			},
			expected: true, // should still create scope
		},
		{
			name:   "multiple fields with empty field",
			fields: []string{"name", "", "title"},
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"test"},
			},
			expected: true,
		},
		{
			name:   "unsupported operator",
			fields: []string{"name"},
			operation: FilterOperation{
				Operator: OperatorGt, // not typically supported for strings
				Values:   []string{"test"},
			},
			expected: true, // function should still return scope but skip unsupported operators
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildStringScope(tt.fields, tt.operation)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// TestBuildBoolScope_EdgeCases tests edge cases in boolean scope building
func TestBuildBoolScope_EdgeCases(t *testing.T) {
	sb := &ScopeBuilder{}

	tests := []struct {
		name      string
		field     string
		operation FilterOperation
		expected  bool // whether scope should be created
	}{
		{
			name:  "unsupported operator",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorLike, // not supported for bool
				Values:   []string{"true"},
			},
			expected: false,
		},
		{
			name:  "empty field name",
			field: "",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"true"},
			},
			expected: true, // should still create scope
		},
		{
			name:  "case variations",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"TRUE"},
			},
			expected: true,
		},
		{
			name:  "numeric string",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"1"},
			},
			expected: true, // should be treated as false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildBoolScope(tt.field, tt.operation)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// TestBuildDateScope_EdgeCases tests edge cases in date scope building
func TestBuildDateScope_EdgeCases(t *testing.T) {
	sb := &ScopeBuilder{}

	tests := []struct {
		name      string
		field     string
		operation FilterOperation
		expected  bool // whether scope should be created
	}{
		{
			name:  "invalid date format",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"not-a-date"},
			},
			expected: false,
		},
		{
			name:  "empty date string",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{""},
			},
			expected: false,
		},
		{
			name:  "between with invalid start date",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"invalid", "2023-12-31"},
			},
			expected: false,
		},
		{
			name:  "between with invalid end date",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"2023-01-01", "invalid"},
			},
			expected: false,
		},
		{
			name:  "unsupported operator",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorLike, // not supported for dates
				Values:   []string{"2023-01-01"},
			},
			expected: false,
		},
		{
			name:  "partial date format",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"2023-01"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildDateScope(tt.field, tt.operation)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// TestBuildEnumScope_EdgeCases tests edge cases in enum scope building
func TestBuildEnumScope_EdgeCases(t *testing.T) {
	sb := &ScopeBuilder{}
	allowedValues := []string{"active", "inactive", "pending"}

	tests := []struct {
		name      string
		field     string
		operation FilterOperation
		expected  bool // whether scope should be created
	}{
		{
			name:  "invalid enum value",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"invalid_status"},
			},
			expected: false,
		},
		{
			name:  "mixed valid and invalid values",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"active", "invalid_status"},
			},
			expected: false,
		},
		{
			name:  "empty allowed values",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"active"},
			},
			expected: false, // will test with empty allowedValues
		},
		{
			name:  "unsupported operator",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorLike, // not supported for enum
				Values:   []string{"active"},
			},
			expected: false,
		},
		{
			name:  "case sensitive enum",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"ACTIVE"}, // different case
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowedVals := allowedValues
			if tt.name == "empty allowed values" {
				allowedVals = []string{}
			}
			scope := sb.buildEnumScope(tt.field, tt.operation, allowedVals)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// TestBuildFilterScope_EdgeCases tests edge cases in filter scope building
func TestBuildFilterScope_EdgeCases(t *testing.T) {
	sb := &ScopeBuilder{}

	tests := []struct {
		name      string
		config    FilterConfig
		values    []string
		expected  bool // whether scope should be created
	}{
		{
			name: "empty field name with no search fields",
			config: FilterConfig{
				Field:        "", // empty field
				SearchFields: []string{}, // no search fields
				Type:         FilterTypeString,
				Operators:    []FilterOperator{OperatorEquals},
			},
			values:   []string{"test"},
			expected: false, // should return nil for empty field
		},
		{
			name: "table name with empty field",
			config: FilterConfig{
				Field:     "",
				TableName: "users",
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals},
			},
			values:   []string{"test"},
			expected: true, // should work with table name
		},
		{
			name: "search fields with table name",
			config: FilterConfig{
				Field:        "name",
				SearchFields: []string{"first_name", "last_name"},
				TableName:    "users",
				Type:         FilterTypeString,
				Operators:    []FilterOperator{OperatorEquals},
			},
			values:   []string{"test"},
			expected: true,
		},
		{
			name: "unknown filter type",
			config: FilterConfig{
				Field:     "custom_field",
				Type:      FilterType("unknown"), // unknown type
				Operators: []FilterOperator{OperatorEquals},
			},
			values:   []string{"test"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildFilterScope(tt.config, tt.values)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// TestBuildMetaScopes_EdgeCases tests edge cases in meta scope building
func TestBuildMetaScopes_EdgeCases(t *testing.T) {
	tests := []struct {
		name               string
		pagination         *Pagination
		expectedScopeCount int
	}{
		{
			name: "zero limit with max limit",
			pagination: &Pagination{
				Limit: 0,
				Options: PaginationOptions{
					DefaultLimit: 10,
					MaxLimit:     100,
				},
			},
			expectedScopeCount: 2, // limit and order (no offset)
		},
		{
			name: "limit exceeds max limit",
			pagination: &Pagination{
				Limit: 200, // exceeds max
				Options: PaginationOptions{
					DefaultLimit: 10,
					MaxLimit:     100,
				},
			},
			expectedScopeCount: 2, // should use default limit
		},
		{
			name: "negative offset",
			pagination: &Pagination{
				Offset: -5,
				Options: PaginationOptions{
					DefaultLimit: 10,
				},
			},
			expectedScopeCount: 2, // no offset scope for negative values
		},
		{
			name: "zero offset",
			pagination: &Pagination{
				Offset: 0,
				Options: PaginationOptions{
					DefaultLimit: 10,
				},
			},
			expectedScopeCount: 2, // no offset scope for zero
		},
		{
			name: "empty order strings",
			pagination: &Pagination{
				Order: "",
				Options: PaginationOptions{
					DefaultLimit: 10,
					DefaultOrder: "", // also empty
				},
			},
			expectedScopeCount: 2, // still creates order scope even if empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &ScopeBuilder{
				pagination: tt.pagination,
				metaScopes: make([]func(*gorm.DB) *gorm.DB, 0),
			}
			sb.buildMetaScopes()
			assert.Equal(t, tt.expectedScopeCount, len(sb.metaScopes))
		})
	}
}

// TestScopeBuilder_DryRun_EdgeCases tests edge cases using dry run
func TestScopeBuilder_DryRun_EdgeCases(t *testing.T) {
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

	t.Run("string scope with multiple fields and OR conditions", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("name", FilterConfig{
			Field:        "name",
			SearchFields: []string{"first_name", "last_name", "nickname"},
			Type:         FilterTypeString,
			Operators:    []FilterOperator{OperatorEquals, OperatorLike},
		})

		conditions := map[string][]string{
			"name": {"like:John"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)
		filterScopes, _ := scopeBuilder.Build()

		var models []TestModel
		query := dryRunDB.Model(&TestModel{})

		// Apply filter scopes
		for _, scope := range filterScopes {
			query = scope(query)
		}

		stmt := query.Find(&models).Statement
		sql := strings.ToUpper(stmt.SQL.String())

		// Should contain OR conditions for multiple fields
		assert.Contains(t, sql, "OR")
		assert.Contains(t, sql, "LIKE")
		assert.Contains(t, sql, "FIRST_NAME")
		assert.Contains(t, sql, "LAST_NAME")
		assert.Contains(t, sql, "NICKNAME")
	})

	t.Run("enum scope with invalid values should not create scope", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			Operators:  []FilterOperator{OperatorEquals, OperatorIn},
			EnumValues: []string{"active", "inactive"},
		})

		conditions := map[string][]string{
			"status": {"invalid_status"}, // invalid enum value
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)
		filterScopes, _ := scopeBuilder.Build()

		// Should not create any filter scopes for invalid enum values
		assert.Equal(t, 0, len(filterScopes))
	})

	t.Run("date scope with timezone handling", func(t *testing.T) {
		filterDef := NewFilterDefinition().AddFilter("created_at", FilterConfig{
			Field:     "created_at",
			Type:      FilterTypeDate,
			Operators: []FilterOperator{OperatorEquals, OperatorBetween},
		})

		conditions := map[string][]string{
			"created_at": {"2023-01-01"},
		}

		pagination := NewPagination(conditions, filterDef, PaginationOptions{})
		scopeBuilder := NewScopeBuilder(pagination)
		filterScopes, _ := scopeBuilder.Build()

		var models []TestModel
		query := dryRunDB.Model(&TestModel{})

		// Apply filter scopes
		for _, scope := range filterScopes {
			query = scope(query)
		}

		stmt := query.Find(&models).Statement
		sql := strings.ToUpper(stmt.SQL.String())

		// Should use BETWEEN for date equals to cover full day
		assert.Contains(t, sql, "BETWEEN")
		assert.Len(t, stmt.Vars, 2) // start and end of day

		// Verify timezone handling (Asia/Jakarta)
		for _, v := range stmt.Vars {
			if timeVal, ok := v.(time.Time); ok {
				assert.Equal(t, "Asia/Jakarta", timeVal.Location().String())
			}
		}
	})
}

// TestScopeBuilder_NilHandling tests nil pointer handling
func TestScopeBuilder_NilHandling(t *testing.T) {
	t.Run("build with nil pagination", func(t *testing.T) {
		sb := &ScopeBuilder{
			pagination: nil,
			scopes:     make([]func(*gorm.DB) *gorm.DB, 0),
			metaScopes: make([]func(*gorm.DB) *gorm.DB, 0),
		}

		// Should not panic
		assert.NotPanics(t, func() {
			filterScopes, metaScopes := sb.Build()
			assert.NotNil(t, filterScopes)
			assert.NotNil(t, metaScopes)
		})
	})

	t.Run("build filter scopes with nil conditions", func(t *testing.T) {
		pagination := &Pagination{
			Conditions: nil, // nil conditions
			FilterDef:  NewFilterDefinition(),
		}
		sb := &ScopeBuilder{
			pagination: pagination,
			scopes:     make([]func(*gorm.DB) *gorm.DB, 0),
			metaScopes: make([]func(*gorm.DB) *gorm.DB, 0),
		}

		assert.NotPanics(t, func() {
			sb.buildFilterScopes()
		})
	})

	t.Run("build filter scopes with nil filter definition", func(t *testing.T) {
		pagination := &Pagination{
			Conditions: map[string][]string{"name": {"test"}},
			FilterDef:  nil, // nil filter definition
		}
		sb := &ScopeBuilder{
			pagination: pagination,
			scopes:     make([]func(*gorm.DB) *gorm.DB, 0),
			metaScopes: make([]func(*gorm.DB) *gorm.DB, 0),
		}

		// Should not panic but also not create any scopes
		assert.NotPanics(t, func() {
			sb.buildFilterScopes()
			assert.Equal(t, 0, len(sb.scopes))
		})
	})
}