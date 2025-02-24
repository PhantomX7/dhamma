package pagination

import (
	"strconv"

	"github.com/PhantomX7/dhamma/utility/scope"
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
	TableName  string // For joined tables
	Operators  []FilterOperator
	EnumValues []string // For enum type validation
}

type FilterDefinition struct {
	configs map[string]FilterConfig
}

func NewFilterDefinition() *FilterDefinition {
	return &FilterDefinition{
		configs: make(map[string]FilterConfig),
	}
}

func (fd *FilterDefinition) AddFilter(field string, config FilterConfig) *FilterDefinition {
	fd.configs[field] = config
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

	limit := 0
	offset := 0
	if len(conditions["limit"]) > 0 {
		var err error
		if limit, err = strconv.Atoi(conditions["limit"][0]); err != nil {
			limit = 0
		}
	}
	if len(conditions["offset"]) > 0 {
		var err error
		if offset, err = strconv.Atoi(conditions["offset"][0]); err != nil {
			offset = 0
		}
	}

	return &Pagination{
		Conditions: conditions,
		FilterDef:  filterDef,
		Options:    options,
		Limit:      limit,
		Offset:     offset,
	}
}

func (p *Pagination) AddCustomScope(scopes ...scope.Scope) {
	p.customScopes = append(p.customScopes, scopes...)
}
