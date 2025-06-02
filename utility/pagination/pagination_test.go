package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Test FilterDefinition
func TestNewFilterDefinition(t *testing.T) {
	fd := NewFilterDefinition()
	assert.NotNil(t, fd)
	assert.NotNil(t, fd.configs)
	assert.NotNil(t, fd.sorts)
	assert.Equal(t, 0, len(fd.configs))
	assert.Equal(t, 0, len(fd.sorts))
}

func TestFilterDefinition_AddFilter(t *testing.T) {
	fd := NewFilterDefinition()
	config := FilterConfig{
		Field:     "name",
		Type:      FilterTypeString,
		Operators: []FilterOperator{OperatorEquals, OperatorLike},
	}

	result := fd.AddFilter("name", config)
	assert.Equal(t, fd, result) // Should return self for chaining
	assert.Equal(t, config, fd.configs["name"])
}

func TestFilterDefinition_AddSort(t *testing.T) {
	fd := NewFilterDefinition()
	config := SortConfig{
		Field:   "created_at",
		Allowed: true,
	}

	result := fd.AddSort("created_at", config)
	assert.Equal(t, fd, result) // Should return self for chaining
	assert.Equal(t, config, fd.sorts["created_at"])
}

// Test Pagination
func TestNewPagination(t *testing.T) {
	tests := []struct {
		name           string
		conditions     map[string][]string
		filterDef      *FilterDefinition
		options        PaginationOptions
		expectedLimit  int
		expectedOffset int
		expectedOrder  string
	}{
		{
			name:           "default values",
			conditions:     map[string][]string{},
			filterDef:      NewFilterDefinition(),
			options:        PaginationOptions{},
			expectedLimit:  20,
			expectedOffset: 0,
			expectedOrder:  "id desc",
		},
		{
			name:       "custom options",
			conditions: map[string][]string{},
			filterDef:  NewFilterDefinition(),
			options: PaginationOptions{
				DefaultLimit: 50,
				MaxLimit:     200,
				DefaultOrder: "name asc",
			},
			expectedLimit:  50,
			expectedOffset: 0,
			expectedOrder:  "name asc",
		},
		{
			name: "with query parameters",
			conditions: map[string][]string{
				"limit":  {"30"},
				"offset": {"10"},
				"sort":   {"created_at desc"},
			},
			filterDef: NewFilterDefinition().AddSort("created_at", SortConfig{Field: "created_at", Allowed: true}),
			options: PaginationOptions{
				DefaultLimit: 20,
				MaxLimit:     100,
				DefaultOrder: "id desc",
			},
			expectedLimit:  30,
			expectedOffset: 10,
			expectedOrder:  "created_at desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination(tt.conditions, tt.filterDef, tt.options)
			assert.Equal(t, tt.expectedLimit, p.Limit)
			assert.Equal(t, tt.expectedOffset, p.Offset)
			assert.Equal(t, tt.expectedOrder, p.Order)
			assert.Equal(t, tt.conditions, p.Conditions)
			assert.Equal(t, tt.filterDef, p.FilterDef)
		})
	}
}

func TestPagination_AddCustomScope(t *testing.T) {
	p := NewPagination(map[string][]string{}, NewFilterDefinition(), PaginationOptions{})
	assert.Equal(t, 0, len(p.customScopes))

	scope1 := func(db *gorm.DB) *gorm.DB { return db.Where("active = ?", true) }
	scope2 := func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }

	p.AddCustomScope(scope1, scope2)
	assert.Equal(t, 2, len(p.customScopes))
}

// Test parsing functions
func TestParseLimit(t *testing.T) {
	tests := []struct {
		name         string
		conditions   map[string][]string
		defaultLimit int
		maxLimit     int
		expected     int
	}{
		{
			name:         "no limit in conditions",
			conditions:   map[string][]string{},
			defaultLimit: 20,
			maxLimit:     100,
			expected:     20,
		},
		{
			name:         "valid limit",
			conditions:   map[string][]string{"limit": {"50"}},
			defaultLimit: 20,
			maxLimit:     100,
			expected:     50,
		},
		{
			name:         "limit exceeds max",
			conditions:   map[string][]string{"limit": {"150"}},
			defaultLimit: 20,
			maxLimit:     100,
			expected:     20,
		},
		{
			name:         "invalid limit",
			conditions:   map[string][]string{"limit": {"invalid"}},
			defaultLimit: 20,
			maxLimit:     100,
			expected:     20,
		},
		{
			name:         "negative limit",
			conditions:   map[string][]string{"limit": {"-10"}},
			defaultLimit: 20,
			maxLimit:     100,
			expected:     20,
		},
		{
			name:         "zero limit",
			conditions:   map[string][]string{"limit": {"0"}},
			defaultLimit: 20,
			maxLimit:     100,
			expected:     20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLimit(tt.conditions, tt.defaultLimit, tt.maxLimit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseOffset(t *testing.T) {
	tests := []struct {
		name       string
		conditions map[string][]string
		expected   int
	}{
		{
			name:       "no offset in conditions",
			conditions: map[string][]string{},
			expected:   0,
		},
		{
			name:       "valid offset",
			conditions: map[string][]string{"offset": {"10"}},
			expected:   10,
		},
		{
			name:       "invalid offset",
			conditions: map[string][]string{"offset": {"invalid"}},
			expected:   0,
		},
		{
			name:       "negative offset",
			conditions: map[string][]string{"offset": {"-5"}},
			expected:   0,
		},
		{
			name:       "zero offset",
			conditions: map[string][]string{"offset": {"0"}},
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOffset(tt.conditions)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateOrder(t *testing.T) {
	filterDef := NewFilterDefinition().
		AddSort("name", SortConfig{Field: "name", Allowed: true}).
		AddSort("created_at", SortConfig{Field: "created_at", Allowed: true}).
		AddSort("disabled_field", SortConfig{Field: "disabled_field", Allowed: false})

	tests := []struct {
		name     string
		order    string
		expected bool
	}{
		{
			name:     "empty order",
			order:    "",
			expected: true,
		},
		{
			name:     "valid single field",
			order:    "name",
			expected: true,
		},
		{
			name:     "valid single field with direction",
			order:    "name desc",
			expected: true,
		},
		{
			name:     "valid multiple fields",
			order:    "name asc, created_at desc",
			expected: true,
		},
		{
			name:     "invalid field",
			order:    "invalid_field",
			expected: false,
		},
		{
			name:     "disabled field",
			order:    "disabled_field",
			expected: false,
		},
		{
			name:     "mixed valid and invalid",
			order:    "name, invalid_field",
			expected: false,
		},
		{
			name:     "empty field part",
			order:    ", name",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateOrder(tt.order, filterDef)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseOrder(t *testing.T) {
	filterDef := NewFilterDefinition().
		AddSort("name", SortConfig{Field: "name", Allowed: true})

	tests := []struct {
		name         string
		conditions   map[string][]string
		defaultOrder string
		expected     string
	}{
		{
			name:         "no sort in conditions",
			conditions:   map[string][]string{},
			defaultOrder: "id desc",
			expected:     "id desc",
		},
		{
			name:         "valid sort",
			conditions:   map[string][]string{"sort": {"name asc"}},
			defaultOrder: "id desc",
			expected:     "name asc",
		},
		{
			name:         "invalid sort",
			conditions:   map[string][]string{"sort": {"invalid_field"}},
			defaultOrder: "id desc",
			expected:     "id desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOrder(tt.conditions, tt.defaultOrder, filterDef)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test constants
func TestConstants(t *testing.T) {
	assert.Equal(t, "limit", QueryKeyLimit)
	assert.Equal(t, "offset", QueryKeyOffset)
	assert.Equal(t, "sort", QueryKeySort)

	// Test FilterType constants
	assert.Equal(t, FilterType("ID"), FilterTypeID)
	assert.Equal(t, FilterType("NUMBER"), FilterTypeNumber)
	assert.Equal(t, FilterType("STRING"), FilterTypeString)
	assert.Equal(t, FilterType("BOOL"), FilterTypeBool)
	assert.Equal(t, FilterType("DATE"), FilterTypeDate)
	assert.Equal(t, FilterType("DATETIME"), FilterTypeDateTime)
	assert.Equal(t, FilterType("ENUM"), FilterTypeEnum)

	// Test FilterOperator constants
	assert.Equal(t, FilterOperator("eq"), OperatorEquals)
	assert.Equal(t, FilterOperator("neq"), OperatorNotEquals)
	assert.Equal(t, FilterOperator("in"), OperatorIn)
	assert.Equal(t, FilterOperator("not_in"), OperatorNotIn)
	assert.Equal(t, FilterOperator("like"), OperatorLike)
	assert.Equal(t, FilterOperator("between"), OperatorBetween)
	assert.Equal(t, FilterOperator("gt"), OperatorGt)
	assert.Equal(t, FilterOperator("gte"), OperatorGte)
	assert.Equal(t, FilterOperator("lt"), OperatorLt)
	assert.Equal(t, FilterOperator("lte"), OperatorLte)
}

// Integration test
func TestPaginationIntegration(t *testing.T) {
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
		AddSort("name", SortConfig{Field: "name", Allowed: true}).
		AddSort("created_at", SortConfig{Field: "created_at", Allowed: true})

	conditions := map[string][]string{
		"limit":  {"25"},
		"offset": {"50"},
		"sort":   {"name asc, created_at desc"},
		"name":   {"like:john"},
		"age":    {"between:20,30"},
	}

	options := PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	}

	p := NewPagination(conditions, filterDef, options)

	assert.Equal(t, 25, p.Limit)
	assert.Equal(t, 50, p.Offset)
	assert.Equal(t, "name asc, created_at desc", p.Order)
	assert.Equal(t, conditions, p.Conditions)
	assert.Equal(t, filterDef, p.FilterDef)
	assert.Equal(t, options, p.Options)
}

// Benchmark tests
func BenchmarkNewPagination(b *testing.B) {
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike},
		})

	conditions := map[string][]string{
		"limit":  {"20"},
		"offset": {"0"},
		"name":   {"like:test"},
	}

	options := PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewPagination(conditions, filterDef, options)
	}
}

func BenchmarkParseLimit(b *testing.B) {
	conditions := map[string][]string{"limit": {"50"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseLimit(conditions, 20, 100)
	}
}

func BenchmarkValidateOrder(b *testing.B) {
	filterDef := NewFilterDefinition().
		AddSort("name", SortConfig{Field: "name", Allowed: true}).
		AddSort("created_at", SortConfig{Field: "created_at", Allowed: true})

	order := "name asc, created_at desc"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateOrder(order, filterDef)
	}
}
