// request_util/types.go
package pagination

import "github.com/PhantomX7/dhamma/utility/scope"

type FilterType string

const (
	FilterTypeID       FilterType = "ID"
	FilterTypeNumber   FilterType = "NUMBER"
	FilterTypeString   FilterType = "STRING"
	FilterTypeBool     FilterType = "BOOL"
	FilterTypeDate     FilterType = "DATE"
	FilterTypeDateTime FilterType = "DATETIME"
	FilterTypeEnum     FilterType = "ENUM"
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

// request_util/filter.go
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

	return &Pagination{
		Conditions: conditions,
		FilterDef:  filterDef,
		Options:    options,
	}
}
