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
	if sb.pagination != nil {
		sb.buildFilterScopes()
		sb.buildMetaScopes()
	}
	return sb.scopes, sb.metaScopes
}

func (sb *ScopeBuilder) buildFilterScopes() {
	if sb.pagination == nil || sb.pagination.Conditions == nil {
		return
	}

	for field, values := range sb.pagination.Conditions {
		if sb.pagination.FilterDef != nil {
			if config, exists := sb.pagination.FilterDef.configs[field]; exists {
				if scope := sb.buildFilterScope(config, values); scope != nil {
					sb.scopes = append(sb.scopes, scope)
				}
			}
		}
	}

	if sb.pagination.customScopes != nil {
		sb.scopes = append(sb.scopes, sb.pagination.customScopes...)
	}
}

func (sb *ScopeBuilder) buildMetaScopes() {
	if sb.pagination == nil {
		return
	}

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

	// Use SearchFields if available, otherwise use Field
	fieldsToSearch := config.SearchFields
	if len(fieldsToSearch) == 0 {
		fieldsToSearch = []string{config.Field}
	}

	// Apply TableName prefix if specified and field doesn't already have a table prefix
	for i, field := range fieldsToSearch {
		if config.TableName != "" && !strings.Contains(field, ".") {
			fieldsToSearch[i] = fmt.Sprintf("%s.%s", config.TableName, field)
		} else {
			// Ensure the field name is not empty if TableName is not used
			if field == "" && len(config.SearchFields) == 0 { // only check if it's the original single Field
				return nil // Or handle error appropriately
			}
		}
	}

	switch config.Type {
	case FilterTypeID, FilterTypeNumber:
		// For numbers, typically we don't search across multiple fields with OR using a single query param.
		// If SearchFields is used with Number type, it will apply the condition to the first field in SearchFields.
		// Consider if this behavior needs adjustment based on specific use cases.
		return sb.buildNumberScope(fieldsToSearch[0], operation)
	case FilterTypeString:
		return sb.buildStringScope(fieldsToSearch, operation)
	case FilterTypeBool:
		return sb.buildBoolScope(fieldsToSearch[0], operation) // Similar to Number, applies to the first field
	case FilterTypeDate, FilterTypeDateTime:
		return sb.buildDateScope(fieldsToSearch[0], operation) // Similar to Number, applies to the first field
	case FilterTypeEnum:
		return sb.buildEnumScope(fieldsToSearch[0], operation, config.EnumValues) // Similar to Number, applies to the first field
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

func (sb *ScopeBuilder) buildStringScope(fields []string, op FilterOperation) func(*gorm.DB) *gorm.DB {
	if len(fields) == 0 {
		return nil
	}

	// If only one field, use simple WHERE clause
	if len(fields) == 1 {
		field := fields[0]
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
		default:
			// Return no-op scope for unsupported operators instead of nil
			return func(db *gorm.DB) *gorm.DB {
				return db
			}
		}
	}

	// For multiple fields, build OR conditions properly
	return func(db *gorm.DB) *gorm.DB {
		var conditions []string
		var args []interface{}

		for _, field := range fields {
			switch op.Operator {
			case OperatorEquals:
				conditions = append(conditions, fmt.Sprintf("%s = ?", field))
				args = append(args, op.Values[0])
			case OperatorNotEquals:
				conditions = append(conditions, fmt.Sprintf("%s != ?", field))
				args = append(args, op.Values[0])
			case OperatorLike:
				conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
				args = append(args, fmt.Sprintf("%%%s%%", op.Values[0]))
			case OperatorIn:
				conditions = append(conditions, fmt.Sprintf("%s IN (?)", field))
				args = append(args, op.Values)
			case OperatorNotIn:
				conditions = append(conditions, fmt.Sprintf("%s NOT IN (?)", field))
				args = append(args, op.Values)
			default:
				// Skip unsupported operator for this field
				continue
			}
		}

		if len(conditions) == 0 {
			return db
		}

		// Join conditions with OR and wrap in parentheses
		orClause := fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
		return db.Where(orClause, args...)
	}
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

func (sb *ScopeBuilder) buildDateScope(field string, op FilterOperation) func(*gorm.DB) *gorm.DB {
	location, err := time.LoadLocation("Asia/Jakarta") // GMT+7
	if err != nil {
		return nil
	}

	switch op.Operator {
	case OperatorEquals:
		t, err := time.ParseInLocation("2006-01-02", op.Values[0], location)
		if err != nil {
			return nil
		}
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(
				fmt.Sprintf("%s BETWEEN ? AND ?", field),
				t,
				t.Add(24*time.Hour).Add(-time.Second),
			)
		}
	case OperatorBetween:
		start, err := time.ParseInLocation("2006-01-02", op.Values[0], location)
		if err != nil {
			return nil
		}
		end, err := time.ParseInLocation("2006-01-02", op.Values[1], location)
		if err != nil {
			return nil
		}

		// Adjust end date to include the entire day (23:59:59.999999999)
		end = end.Add(24 * time.Hour).Add(-time.Nanosecond)

		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), start, end)
		}
	case OperatorGte:
		t, err := time.ParseInLocation("2006-01-02", op.Values[0], location)
		if err != nil {
			return nil
		}
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s >= ?", field), t)
		}
	case OperatorLte:
		t, err := time.ParseInLocation("2006-01-02", op.Values[0], location)
		if err != nil {
			return nil
		}
		// Adjust to include the entire day (23:59:59.999999999)
		t = t.Add(24 * time.Hour).Add(-time.Nanosecond)
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s <= ?", field), t)
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
