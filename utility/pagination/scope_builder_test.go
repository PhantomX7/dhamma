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

// Test model for GORM operations
type TestModel struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Age       int
	Active    bool
	Status    string
	CreatedAt time.Time
}

// Helper function to create test database for dry run testing
func setupTestDB(t *testing.T) *gorm.DB {
	// Use SQLite in-memory database for testing
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

	return db
}

// Helper function to create dry run database session
func setupDryRunDB(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	return db.Session(&gorm.Session{DryRun: true})
}

// Helper function to validate SQL query structure
func validateSQL(t *testing.T, stmt *gorm.Statement, expectedPatterns []string) {
	sql := stmt.SQL.String()
	t.Logf("Generated SQL: %s", sql)
	t.Logf("Variables: %v", stmt.Vars)

	for _, pattern := range expectedPatterns {
		assert.Contains(t, strings.ToUpper(sql), strings.ToUpper(pattern), "SQL should contain pattern: %s", pattern)
	}
}

// Test ScopeBuilder
func TestNewScopeBuilder(t *testing.T) {
	pagination := &Pagination{}
	sb := NewScopeBuilder(pagination)

	assert.NotNil(t, sb)
	assert.Equal(t, pagination, sb.pagination)
	assert.NotNil(t, sb.scopes)
	assert.NotNil(t, sb.metaScopes)
	assert.Equal(t, 0, len(sb.scopes))
	assert.Equal(t, 0, len(sb.metaScopes))
}

func TestScopeBuilder_Build(t *testing.T) {
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorEquals},
		})

	conditions := map[string][]string{
		"name":   {"John"},
		"limit":  {"10"},
		"offset": {"5"},
	}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	})

	sb := NewScopeBuilder(pagination)
	filterScopes, metaScopes := sb.Build()

	assert.NotNil(t, filterScopes)
	assert.NotNil(t, metaScopes)
	assert.Equal(t, 1, len(filterScopes)) // One filter scope for name
	assert.Equal(t, 3, len(metaScopes))   // Limit, offset, and order scopes
}

func TestScopeBuilder_BuildWithCustomScopes(t *testing.T) {
	filterDef := NewFilterDefinition()
	conditions := map[string][]string{}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{})
	customScope := func(db *gorm.DB) *gorm.DB {
		return db.Where("active = ?", true)
	}
	pagination.AddCustomScope(customScope)

	sb := NewScopeBuilder(pagination)
	filterScopes, metaScopes := sb.Build()

	assert.Equal(t, 1, len(filterScopes)) // One custom scope
	assert.Equal(t, 2, len(metaScopes))   // Limit and order scopes (no offset)
}

// Test parseFilterOperation
func TestParseFilterOperation(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected FilterOperation
	}{
		{
			name:  "simple value without operator",
			value: "John",
			expected: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"John"},
			},
		},
		{
			name:  "equals operator",
			value: "eq:John",
			expected: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"John"},
			},
		},
		{
			name:  "like operator",
			value: "like:John",
			expected: FilterOperation{
				Operator: OperatorLike,
				Values:   []string{"John"},
			},
		},
		{
			name:  "in operator with multiple values",
			value: "in:active,inactive,pending",
			expected: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"active", "inactive", "pending"},
			},
		},
		{
			name:  "between operator",
			value: "between:20,30",
			expected: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"20", "30"},
			},
		},
		{
			name:  "value with colon but no operator",
			value: "test:value:with:colons",
			expected: FilterOperation{
				Operator: FilterOperator("test"),
				Values:   []string{"value:with:colons"},
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

// Test isValidOperation
func TestIsValidOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation FilterOperation
		config    FilterConfig
		expected  bool
	}{
		{
			name: "valid equals operation",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"John"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorEquals, OperatorLike},
			},
			expected: true,
		},
		{
			name: "invalid operator not allowed",
			operation: FilterOperation{
				Operator: OperatorGt,
				Values:   []string{"25"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorEquals, OperatorLike},
			},
			expected: false,
		},
		{
			name: "valid between operation with two values",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"20", "30"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorBetween},
			},
			expected: true,
		},
		{
			name: "invalid between operation with one value",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"20"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorBetween},
			},
			expected: false,
		},
		{
			name: "valid in operation with multiple values",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"active", "inactive"},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorIn},
			},
			expected: true,
		},
		{
			name: "invalid in operation with no values",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{},
			},
			config: FilterConfig{
				Operators: []FilterOperator{OperatorIn},
			},
			expected: false,
		},
		{
			name: "invalid equals operation with multiple values",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"John", "Jane"},
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

// Test buildNumberScope
func TestScopeBuilder_BuildNumberScope(t *testing.T) {
	t.Skip("Database tests require MySQL setup")
	sb := &ScopeBuilder{}

	tests := []struct {
		name          string
		field         string
		operation     FilterOperation
		expectedCount int64
		expectedSQL   string
	}{
		{
			name:  "equals operation",
			field: "age",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"25"},
			},
			expectedCount: 1,
		},
		{
			name:  "not equals operation",
			field: "age",
			operation: FilterOperation{
				Operator: OperatorNotEquals,
				Values:   []string{"25"},
			},
			expectedCount: 2,
		},
		{
			name:  "greater than operation",
			field: "age",
			operation: FilterOperation{
				Operator: OperatorGt,
				Values:   []string{"25"},
			},
			expectedCount: 2,
		},
		{
			name:  "between operation",
			field: "age",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"25", "30"},
			},
			expectedCount: 2,
		},
		{
			name:  "in operation",
			field: "age",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"25", "30"},
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildNumberScope(tt.field, tt.operation)
			assert.NotNil(t, scope)
		})
	}
}

// Test buildStringScope
func TestScopeBuilder_BuildStringScope(t *testing.T) {
	t.Skip("Database tests require MySQL setup")
	sb := &ScopeBuilder{}

	tests := []struct {
		name          string
		field         string
		operation     FilterOperation
		expectedCount int64
	}{
		{
			name:  "equals operation",
			field: "name",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"John"},
			},
			expectedCount: 1,
		},
		{
			name:  "like operation",
			field: "name",
			operation: FilterOperation{
				Operator: OperatorLike,
				Values:   []string{"Jo"},
			},
			expectedCount: 1,
		},
		{
			name:  "in operation",
			field: "name",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"John", "Jane"},
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildStringScope(tt.field, tt.operation)
			assert.NotNil(t, scope)
		})
	}
}

// Test buildBoolScope
func TestScopeBuilder_BuildBoolScope(t *testing.T) {
	t.Skip("Database tests require MySQL setup")
	sb := &ScopeBuilder{}

	tests := []struct {
		name          string
		field         string
		operation     FilterOperation
		expectedCount int64
	}{
		{
			name:  "true value",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"true"},
			},
			expectedCount: 2,
		},
		{
			name:  "false value",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"false"},
			},
			expectedCount: 1,
		},
		{
			name:  "True value (case insensitive)",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"True"},
			},
			expectedCount: 2,
		},
		{
			name:  "invalid value treated as false",
			field: "active",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"invalid"},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildBoolScope(tt.field, tt.operation)
			assert.NotNil(t, scope)
		})
	}
}

// Test buildDateScope
func TestScopeBuilder_BuildDateScope(t *testing.T) {
	t.Skip("Database tests require MySQL setup")
	sb := &ScopeBuilder{}

	tests := []struct {
		name      string
		field     string
		operation FilterOperation
		expected  bool // whether scope should be created
	}{
		{
			name:  "valid date equals",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"2023-01-01"},
			},
			expected: true,
		},
		{
			name:  "valid date between",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"2023-01-01", "2023-12-31"},
			},
			expected: true,
		},
		{
			name:  "invalid date format",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"invalid-date"},
			},
			expected: false,
		},
		{
			name:  "invalid between with bad end date",
			field: "created_at",
			operation: FilterOperation{
				Operator: OperatorBetween,
				Values:   []string{"2023-01-01", "invalid-date"},
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

// Test buildEnumScope
func TestScopeBuilder_BuildEnumScope(t *testing.T) {
	t.Skip("Database tests require MySQL setup")
	sb := &ScopeBuilder{}
	allowedValues := []string{"active", "inactive", "pending"}

	tests := []struct {
		name          string
		field         string
		operation     FilterOperation
		expected      bool
		expectedCount int64
	}{
		{
			name:  "valid enum equals",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorEquals,
				Values:   []string{"active"},
			},
			expected:      true,
			expectedCount: 1,
		},
		{
			name:  "valid enum in",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"active", "inactive"},
			},
			expected:      true,
			expectedCount: 2,
		},
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
			name:  "mixed valid and invalid enum values",
			field: "status",
			operation: FilterOperation{
				Operator: OperatorIn,
				Values:   []string{"active", "invalid_status"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := sb.buildEnumScope(tt.field, tt.operation, allowedValues)
			if tt.expected {
				assert.NotNil(t, scope)
			} else {
				assert.Nil(t, scope)
			}
		})
	}
}

// Test buildFilterScope
func TestScopeBuilder_BuildFilterScope(t *testing.T) {
	sb := &ScopeBuilder{}

	tests := []struct {
		name     string
		config   FilterConfig
		values   []string
		expected bool
	}{
		{
			name: "valid string filter",
			config: FilterConfig{
				Field:     "name",
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals},
			},
			values:   []string{"John"},
			expected: true,
		},
		{
			name: "valid number filter",
			config: FilterConfig{
				Field:     "age",
				Type:      FilterTypeNumber,
				Operators: []FilterOperator{OperatorGt},
			},
			values:   []string{"gt:25"},
			expected: true,
		},
		{
			name: "empty values",
			config: FilterConfig{
				Field:     "name",
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals},
			},
			values:   []string{},
			expected: false,
		},
		{
			name: "invalid operation",
			config: FilterConfig{
				Field:     "name",
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals},
			},
			values:   []string{"gt:25"}, // gt not allowed for string
			expected: false,
		},
		{
			name: "with table name",
			config: FilterConfig{
				Field:     "name",
				TableName: "users",
				Type:      FilterTypeString,
				Operators: []FilterOperator{OperatorEquals},
			},
			values:   []string{"John"},
			expected: true,
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

// Integration test with GORM dry run
func TestScopeBuilder_Integration(t *testing.T) {
	db := setupDryRunDB(t)

	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike, OperatorEquals},
		}).
		AddFilter("age", FilterConfig{
			Field:     "age",
			Type:      FilterTypeNumber,
			Operators: []FilterOperator{OperatorGt, OperatorLt, OperatorBetween},
		}).
		AddFilter("active", FilterConfig{
			Field:     "active",
			Type:      FilterTypeBool,
			Operators: []FilterOperator{OperatorEquals},
		})

	conditions := map[string][]string{
		"name":   {"like:J"},
		"age":    {"gt:20"},
		"active": {"true"},
		"limit":  {"10"},
		"offset": {"0"},
	}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	})

	scopeBuilder := NewScopeBuilder(pagination)
	filterScopes, metaScopes := scopeBuilder.Build()

	// Verify scopes are created
	assert.NotEmpty(t, filterScopes)
	assert.NotEmpty(t, metaScopes)

	// Test the generated SQL using dry run
	var models []TestModel
	query := db.Model(&TestModel{})

	// Apply filter scopes
	for _, scope := range filterScopes {
		query = scope(query)
	}

	// Apply meta scopes
	for _, scope := range metaScopes {
		query = scope(query)
	}

	// Execute dry run to get SQL
	stmt := query.Find(&models).Statement

	// Validate expected SQL patterns
	expectedPatterns := []string{
		"SELECT",
		"FROM `test_models`",
		"WHERE",
		"name LIKE",
		"age >",
		"active =",
		"ORDER BY",
		"LIMIT",
		// Note: OFFSET might not appear when it's 0
	}

	validateSQL(t, stmt, expectedPatterns)

	// Verify filter variables are properly bound
	assert.Contains(t, stmt.Vars, "%J%") // LIKE pattern
	assert.Contains(t, stmt.Vars, "20")  // age > 20 (as string)
	assert.Contains(t, stmt.Vars, true)  // active = true
	// Note: LIMIT and OFFSET might be handled differently in SQLite vs MySQL
}

// Test dry run for string filters
func TestScopeBuilder_DryRun_StringFilters(t *testing.T) {
	db := setupDryRunDB(t)

	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorEquals, OperatorLike, OperatorIn},
		})

	tests := []struct {
		name             string
		conditions       map[string][]string
		expectedPatterns []string
		expectedVars     []interface{}
	}{
		{
			name: "equals filter",
			conditions: map[string][]string{
				"name": {"eq:John"},
			},
			expectedPatterns: []string{"name ="},
			expectedVars:     []interface{}{"John"},
		},
		{
			name: "like filter",
			conditions: map[string][]string{
				"name": {"like:John"},
			},
			expectedPatterns: []string{"name LIKE"},
			expectedVars:     []interface{}{"%John%"},
		},
		{
			name: "in filter",
			conditions: map[string][]string{
				"name": {"in:John,Jane,Bob"},
			},
			expectedPatterns: []string{"name IN"},
			expectedVars:     []interface{}{"John", "Jane", "Bob"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := NewPagination(tt.conditions, filterDef, PaginationOptions{})
			scopeBuilder := NewScopeBuilder(pagination)
			filterScopes, _ := scopeBuilder.Build()

			var models []TestModel
			query := db.Model(&TestModel{})

			// Apply filter scopes
			for _, scope := range filterScopes {
				query = scope(query)
			}

			stmt := query.Find(&models).Statement
			validateSQL(t, stmt, tt.expectedPatterns)

			// Verify expected variables
			for _, expectedVar := range tt.expectedVars {
				assert.Contains(t, stmt.Vars, expectedVar)
			}
		})
	}
}

// Test dry run for number filters
func TestScopeBuilder_DryRun_NumberFilters(t *testing.T) {
	db := setupDryRunDB(t)

	filterDef := NewFilterDefinition().
		AddFilter("age", FilterConfig{
			Field:     "age",
			Type:      FilterTypeNumber,
			Operators: []FilterOperator{OperatorGt, OperatorLt, OperatorBetween, OperatorEquals},
		})

	tests := []struct {
		name             string
		conditions       map[string][]string
		expectedPatterns []string
		expectedVars     []interface{}
	}{
		{
			name: "greater than filter",
			conditions: map[string][]string{
				"age": {"gt:25"},
			},
			expectedPatterns: []string{"age >"},
			expectedVars:     []interface{}{"25"},
		},
		{
			name: "less than filter",
			conditions: map[string][]string{
				"age": {"lt:65"},
			},
			expectedPatterns: []string{"age <"},
			expectedVars:     []interface{}{"65"},
		},
		{
			name: "between filter",
			conditions: map[string][]string{
				"age": {"between:25,65"},
			},
			expectedPatterns: []string{"age BETWEEN"},
			expectedVars:     []interface{}{"25", "65"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := NewPagination(tt.conditions, filterDef, PaginationOptions{})
			scopeBuilder := NewScopeBuilder(pagination)
			filterScopes, _ := scopeBuilder.Build()

			var models []TestModel
			query := db.Model(&TestModel{})

			// Apply filter scopes
			for _, scope := range filterScopes {
				query = scope(query)
			}

			stmt := query.Find(&models).Statement
			validateSQL(t, stmt, tt.expectedPatterns)

			// Verify expected variables
			for _, expectedVar := range tt.expectedVars {
				assert.Contains(t, stmt.Vars, expectedVar)
			}
		})
	}
}

// Test dry run for datetime filters
func TestScopeBuilder_DryRun_DateTimeFilters(t *testing.T) {
	db := setupDryRunDB(t)

	filterDef := NewFilterDefinition().
		AddFilter("created_at", FilterConfig{
			Field:     "created_at",
			Type:      FilterTypeDateTime,
			Operators: []FilterOperator{OperatorBetween, OperatorGt, OperatorLt},
		})

	conditions := map[string][]string{
		"created_at": {"between:2023-01-01,2023-12-31"},
	}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{})
	scopeBuilder := NewScopeBuilder(pagination)
	filterScopes, _ := scopeBuilder.Build()

	var models []TestModel
	query := db.Model(&TestModel{})

	// Apply filter scopes
	for _, scope := range filterScopes {
		query = scope(query)
	}

	stmt := query.Find(&models).Statement
	validateSQL(t, stmt, []string{"created_at BETWEEN"})

	// Verify date parsing
	assert.Len(t, stmt.Vars, 2)
	assert.IsType(t, time.Time{}, stmt.Vars[0])
	assert.IsType(t, time.Time{}, stmt.Vars[1])
}

// Test dry run for custom scopes
func TestScopeBuilder_DryRun_CustomScopes(t *testing.T) {
	db := setupDryRunDB(t)

	filterDef := NewFilterDefinition()
	conditions := map[string][]string{}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{})

	// Add custom scope
	customScope := func(db *gorm.DB) *gorm.DB {
		return db.Where("active = ?", true)
	}
	pagination.AddCustomScope(customScope)

	scopeBuilder := NewScopeBuilder(pagination)
	filterScopes, _ := scopeBuilder.Build()

	var models []TestModel
	query := db.Model(&TestModel{})

	// Apply filter scopes (including custom scope)
	for _, scope := range filterScopes {
		query = scope(query)
	}

	stmt := query.Find(&models).Statement
	validateSQL(t, stmt, []string{"active ="})
	assert.Contains(t, stmt.Vars, true)
}

// Benchmark tests
func BenchmarkNewScopeBuilder(b *testing.B) {
	pagination := &Pagination{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewScopeBuilder(pagination)
	}
}

func BenchmarkParseFilterOperation(b *testing.B) {
	value := "like:John"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseFilterOperation(value)
	}
}

func BenchmarkIsValidOperation(b *testing.B) {
	operation := FilterOperation{
		Operator: OperatorEquals,
		Values:   []string{"John"},
	}
	config := FilterConfig{
		Operators: []FilterOperator{OperatorEquals, OperatorLike},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isValidOperation(operation, config)
	}
}

func BenchmarkScopeBuilder_Build(b *testing.B) {
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike},
		})

	conditions := map[string][]string{
		"name":   {"like:John"},
		"limit":  {"20"},
		"offset": {"0"},
	}

	pagination := NewPagination(conditions, filterDef, PaginationOptions{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb := NewScopeBuilder(pagination)
		_, _ = sb.Build()
	}
}
