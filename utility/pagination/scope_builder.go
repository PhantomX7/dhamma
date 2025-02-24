// pagination/scope_builder.go
package pagination

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ScopeBuilder struct {
	pagination *Pagination
	scopes     []func(*gorm.DB) *gorm.DB
	metaScopes []func(*gorm.DB) *gorm.DB
}

func NewScopeBuilder(pagination *Pagination) *ScopeBuilder {
	return &ScopeBuilder{
		pagination: pagination,
		scopes:     make([]func(*gorm.DB) *gorm.DB, 0),
		metaScopes: make([]func(*gorm.DB) *gorm.DB, 0),
	}
}

func (sb *ScopeBuilder) Build() ([]func(*gorm.DB) *gorm.DB, []func(*gorm.DB) *gorm.DB) {
	sb.buildFilterScopes()
	sb.buildMetaScopes()
	return sb.scopes, sb.metaScopes
}

func (sb *ScopeBuilder) buildFilterScopes() {
	for field, values := range sb.pagination.Conditions {
		if config, exists := sb.pagination.FilterDef.configs[field]; exists {
			if scope := sb.buildFilterScope(config, values); scope != nil {
				sb.scopes = append(sb.scopes, scope)
			}
		}
	}

	sb.scopes = append(sb.scopes, sb.pagination.customScopes...)
}

func (sb *ScopeBuilder) buildMetaScopes() {
	// Build limit scope
	limit := sb.pagination.Options.DefaultLimit
	if sb.pagination.Limit > 0 && sb.pagination.Limit <= sb.pagination.Options.MaxLimit {
		limit = sb.pagination.Limit
	}
	sb.metaScopes = append(sb.metaScopes, func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})

	// Build offset scope
	if sb.pagination.Offset > 0 {
		sb.metaScopes = append(sb.metaScopes, func(db *gorm.DB) *gorm.DB {
			return db.Offset(sb.pagination.Offset)
		})
	}

	// Build order scope
	order := sb.pagination.Options.DefaultOrder
	if sb.pagination.Order != "" {
		order = sb.pagination.Order
	}
	sb.metaScopes = append(sb.metaScopes, func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	})
}

func (sb *ScopeBuilder) buildFilterScope(config FilterConfig, values []string) func(*gorm.DB) *gorm.DB {
	if len(values) == 0 {
		return nil
	}

	operation := parseFilterOperation(values[0])
	if !isValidOperation(operation, config) {
		return nil
	}

	fieldName := config.Field
	if config.TableName != "" {
		fieldName = fmt.Sprintf("%s.%s", config.TableName, config.Field)
	}

	switch config.Type {
	case FilterTypeID, FilterTypeNumber:
		return sb.buildNumberScope(fieldName, operation)
	case FilterTypeString:
		return sb.buildStringScope(fieldName, operation)
	case FilterTypeBool:
		return sb.buildBoolScope(fieldName, operation)
	case FilterTypeDate, FilterTypeDateTime:
		return sb.buildDateScope(fieldName, operation, config.Type)
	case FilterTypeEnum:
		return sb.buildEnumScope(fieldName, operation, config.EnumValues)
	}

	return nil
}

type FilterOperation struct {
	Operator FilterOperator
	Values   []string
}

func parseFilterOperation(value string) FilterOperation {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return FilterOperation{
			Operator: OperatorEquals,
			Values:   []string{value},
		}
	}

	return FilterOperation{
		Operator: FilterOperator(parts[0]),
		Values:   strings.Split(parts[1], ","),
	}
}

func isValidOperation(operation FilterOperation, config FilterConfig) bool {
	// Check if operator is allowed
	operatorAllowed := false
	for _, allowed := range config.Operators {
		if operation.Operator == allowed {
			operatorAllowed = true
			break
		}
	}
	if !operatorAllowed {
		return false
	}

	// Validate number of values
	switch operation.Operator {
	case OperatorBetween:
		if len(operation.Values) != 2 {
			return false
		}
	case OperatorIn, OperatorNotIn:
		if len(operation.Values) == 0 {
			return false
		}
	default:
		if len(operation.Values) != 1 {
			return false
		}
	}

	return true
}

func (sb *ScopeBuilder) buildNumberScope(field string, op FilterOperation) func(*gorm.DB) *gorm.DB {
	switch op.Operator {
	case OperatorEquals:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", field), op.Values[0])
		}
	case OperatorNotEquals:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s != ?", field), op.Values[0])
		}
	case OperatorIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s IN (?)", field), op.Values)
		}
	case OperatorNotIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s NOT IN (?)", field), op.Values)
		}
	case OperatorBetween:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), op.Values[0], op.Values[1])
		}
	case OperatorGt:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s > ?", field), op.Values[0])
		}
	case OperatorGte:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s >= ?", field), op.Values[0])
		}
	case OperatorLt:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s < ?", field), op.Values[0])
		}
	case OperatorLte:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s <= ?", field), op.Values[0])
		}
	}
	return nil
}

func (sb *ScopeBuilder) buildStringScope(field string, op FilterOperation) func(*gorm.DB) *gorm.DB {
	switch op.Operator {
	case OperatorEquals:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", field), op.Values[0])
		}
	case OperatorNotEquals:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s != ?", field), op.Values[0])
		}
	case OperatorLike:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%s%%", op.Values[0]))
		}
	case OperatorIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s IN (?)", field), op.Values)
		}
	case OperatorNotIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s NOT IN (?)", field), op.Values)
		}
	}
	return nil
}

func (sb *ScopeBuilder) buildBoolScope(field string, op FilterOperation) func(*gorm.DB) *gorm.DB {
	switch op.Operator {
	case OperatorEquals:
		value := strings.ToLower(op.Values[0]) == "true"
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}
	return nil
}

func (sb *ScopeBuilder) buildDateScope(field string, op FilterOperation, filterType FilterType) func(*gorm.DB) *gorm.DB {
	switch op.Operator {
	case OperatorEquals:
		t, err := time.Parse("2006-01-02", op.Values[0])
		if err != nil {
			return nil
		}
		return func(db *gorm.DB) *gorm.DB {
			if filterType == FilterTypeDateTime {
				return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field),
					t, t.Add(24*time.Hour).Add(-time.Second))
			}
			return db.Where(fmt.Sprintf("%s = ?", field), t)
		}
	case OperatorBetween:
		start, err := time.Parse("2006-01-02", op.Values[0])
		if err != nil {
			return nil
		}
		end, err := time.Parse("2006-01-02", op.Values[1])
		if err != nil {
			return nil
		}
		if filterType == FilterTypeDateTime {
			end = end.Add(24 * time.Hour).Add(-time.Second)
		}
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), start, end)
		}
	}
	return nil
}

func (sb *ScopeBuilder) buildEnumScope(field string, op FilterOperation, allowedValues []string) func(*gorm.DB) *gorm.DB {
	// Validate enum values
	for _, val := range op.Values {
		valid := false
		for _, allowedVal := range allowedValues {
			if val == allowedVal {
				valid = true
				break
			}
		}
		if !valid {
			return nil
		}
	}

	switch op.Operator {
	case OperatorEquals:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", field), op.Values[0])
		}
	case OperatorIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s IN (?)", field), op.Values)
		}
	}
	return nil
}
