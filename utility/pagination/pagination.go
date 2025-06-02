package pagination

import (
	"strconv"
	"strings"

	"github.com/PhantomX7/dhamma/utility/scope"
)

const (
	QueryKeyLimit  = "limit"
	QueryKeyOffset = "offset"
	QueryKeySort   = "sort"
)

type FilterType string

// FilterType is a type of filter for query
const (
	FilterTypeID       FilterType = "ID"       // allowed operator: eq, neq, in, not_in, between, gt, gte, lt, lte
	FilterTypeNumber   FilterType = "NUMBER"   // allowed operator: eq, neq, in, not_in, between, gt, gte, lt, lte
	FilterTypeString   FilterType = "STRING"   // allowed operator: eq, neq, in, not_in, like
	FilterTypeBool     FilterType = "BOOL"     // allowed operator: eq
	FilterTypeDate     FilterType = "DATE"     // allowed operator: eq, between
	FilterTypeDateTime FilterType = "DATETIME" // allowed operator: eq, between
	FilterTypeEnum     FilterType = "ENUM"     // allowed operator: eq, in
)

type FilterOperator string

const (
	OperatorEquals    FilterOperator = "eq"
	OperatorNotEquals FilterOperator = "neq"
	OperatorIn        FilterOperator = "in"
	OperatorNotIn     FilterOperator = "not_in"
	OperatorLike      FilterOperator = "like"
	OperatorBetween   FilterOperator = "between"
	OperatorGt        FilterOperator = "gt"
	OperatorGte       FilterOperator = "gte"
	OperatorLt        FilterOperator = "lt"
	OperatorLte       FilterOperator = "lte"
)

type FilterConfig struct {
	Field      string
	Type       FilterType
	TableName  string // For joined tables, use the struct name for Joined field, otherwise use the plural name for own field
	Operators  []FilterOperator
	EnumValues []string // For enum type validation
}

type SortConfig struct {
	Field     string
	TableName string
	Allowed   bool
}

type FilterDefinition struct {
	configs map[string]FilterConfig
	sorts   map[string]SortConfig
}

func NewFilterDefinition() *FilterDefinition {
	return &FilterDefinition{
		configs: make(map[string]FilterConfig),
		sorts:   make(map[string]SortConfig),
	}
}

func (fd *FilterDefinition) AddFilter(field string, config FilterConfig) *FilterDefinition {
	fd.configs[field] = config
	return fd
}

func (fd *FilterDefinition) AddSort(field string, config SortConfig) *FilterDefinition {
	fd.sorts[field] = config
	return fd
}

// request_util/pagination.go
type PaginationOptions struct {
	DefaultLimit int
	MaxLimit     int
	DefaultOrder string
}

type Pagination struct {
	Limit        int
	Offset       int
	Order        string
	Conditions   map[string][]string
	FilterDef    *FilterDefinition
	Options      PaginationOptions
	customScopes []scope.Scope
	Preloads     []string
	Filters      map[string]interface{}
}

func NewPagination(conditions map[string][]string, filterDef *FilterDefinition, options PaginationOptions) *Pagination {
	if options.DefaultLimit == 0 {
		options.DefaultLimit = 20
	}
	if options.MaxLimit == 0 {
		options.MaxLimit = 100
	}
	if options.DefaultOrder == "" {
		options.DefaultOrder = "id desc"
	}

	return &Pagination{
		Conditions: conditions,
		FilterDef:  filterDef,
		Options:    options,
		Limit:      parseLimit(conditions, options.DefaultLimit, options.MaxLimit),
		Offset:     parseOffset(conditions),
		Order:      parseOrder(conditions, options.DefaultOrder, filterDef),
	}
}

func (p *Pagination) AddCustomScope(scopes ...scope.Scope) {
	p.customScopes = append(p.customScopes, scopes...)
}

func parseLimit(conditions map[string][]string, defaultLimit, maxLimit int) int {
	if limitStr, exists := conditions[QueryKeyLimit]; exists && len(limitStr) > 0 {
		if parsedLimit, err := strconv.Atoi(limitStr[0]); err == nil {
			if parsedLimit > 0 && parsedLimit <= maxLimit {
				return parsedLimit
			}
		}
	}
	return defaultLimit
}

func parseOffset(conditions map[string][]string) int {
	if offsetStr, exists := conditions[QueryKeyOffset]; exists && len(offsetStr) > 0 {
		if parsedOffset, err := strconv.Atoi(offsetStr[0]); err == nil && parsedOffset > 0 {
			return parsedOffset
		}
	}
	return 0
}

func validateOrder(order string, filterDef *FilterDefinition) bool {
	if order == "" {
		return true
	}

	for _, part := range strings.Split(order, ",") {
		fieldParts := strings.Fields(strings.TrimSpace(part))
		if len(fieldParts) == 0 {
			return false
		}

		field := fieldParts[0]
		if sortConfig, exists := filterDef.sorts[field]; !exists || !sortConfig.Allowed {
			return false
		}
	}
	return true
}

func parseOrder(conditions map[string][]string, defaultOrder string, filterDef *FilterDefinition) string {
	if orderStr, exists := conditions[QueryKeySort]; exists && len(orderStr) > 0 {
		order := orderStr[0]
		if validateOrder(order, filterDef) {
			return order
		}
	}
	return defaultOrder
}
