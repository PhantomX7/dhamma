package pagination

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ScopeBuilder struct {
	pagination *Pagination
	scopes     []func(*gorm.DB) *gorm.DB
	metaScopes []func(*gorm.DB) *gorm.DB // Separate meta scopes (limit, offset, order)
}

func NewScopeBuilder(pagination *Pagination) *ScopeBuilder {
	return &ScopeBuilder{
		pagination: pagination,
		scopes:     make([]func(*gorm.DB) *gorm.DB, 0),
		metaScopes: make([]func(*gorm.DB) *gorm.DB, 0),
	}
}

func (sb *ScopeBuilder) Build() (filterScopes []func(*gorm.DB) *gorm.DB, metaScopes []func(*gorm.DB) *gorm.DB) {
	sb.buildPaginationScopes()
	sb.buildFilterScopes()
	return sb.scopes, sb.metaScopes
}

func (sb *ScopeBuilder) buildPaginationScopes() {
	// Limit scope
	limit := sb.pagination.Options.DefaultLimit
	if len(sb.pagination.Conditions["limit"]) > 0 {
		if parsedLimit, err := strconv.Atoi(sb.pagination.Conditions["limit"][0]); err == nil {
			if parsedLimit > 0 && parsedLimit <= sb.pagination.Options.MaxLimit {
				limit = parsedLimit
			}
		}
	}
	sb.metaScopes = append(sb.metaScopes, func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})

	// Offset scope
	if len(sb.pagination.Conditions["offset"]) > 0 {
		if offset, err := strconv.Atoi(sb.pagination.Conditions["offset"][0]); err == nil && offset > 0 {
			sb.metaScopes = append(sb.metaScopes, func(db *gorm.DB) *gorm.DB {
				return db.Offset(offset)
			})
		}
	}

	// Order scope
	order := sb.pagination.Options.DefaultOrder
	if len(sb.pagination.Conditions["sort"]) > 0 {
		order = sb.pagination.Conditions["sort"][0]
	}
	sb.metaScopes = append(sb.metaScopes, func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	})
}

func (sb *ScopeBuilder) buildFilterScopes() {
	for field, config := range sb.pagination.FilterDef.configs {
		if values, exists := sb.pagination.Conditions[field]; exists && len(values) > 0 {
			scope := sb.buildFilterScope(config, values)
			if scope != nil {
				sb.scopes = append(sb.scopes, scope)
			}
		}
	}

	sb.scopes = append(sb.scopes, sb.pagination.customScopes...)
}

func (sb *ScopeBuilder) buildFilterScope(config FilterConfig, values []string) func(*gorm.DB) *gorm.DB {
	fieldName := config.Field
	if config.TableName != "" {
		fieldName = fmt.Sprintf("%s.%s", config.TableName, config.Field)
	}

	switch config.Type {
	case FilterTypeID, FilterTypeNumber:
		return sb.buildNumberScope(fieldName, values)
	case FilterTypeString:
		return sb.buildStringScope(fieldName, values)
	case FilterTypeBool:
		return sb.buildBoolScope(fieldName, values)
	case FilterTypeDate, FilterTypeDateTime:
		return sb.buildDateScope(fieldName, values, config.Type)
	case FilterTypeEnum:
		return sb.buildEnumScope(fieldName, values, config.EnumValues)
	}

	return nil
}

func (sb *ScopeBuilder) buildNumberScope(field string, values []string) func(*gorm.DB) *gorm.DB {
	if len(values) == 0 {
		return nil
	}

	// Handle between case
	if strings.Contains(values[0], ",") {
		parts := strings.Split(values[0], ",")
		if len(parts) == 2 {
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), parts[0], parts[1])
			}
		}
	}

	// Handle IN case
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN (?)", field), values)
	}
}

func (sb *ScopeBuilder) buildStringScope(field string, values []string) func(*gorm.DB) *gorm.DB {
	if len(values) == 0 {
		return nil
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%s%%", values[0]))
	}
}

func (sb *ScopeBuilder) buildBoolScope(field string, values []string) func(*gorm.DB) *gorm.DB {
	if len(values) == 0 {
		return nil
	}

	value := strings.ToLower(values[0]) == "true"
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

func (sb *ScopeBuilder) buildDateScope(field string, values []string, filterType FilterType) func(*gorm.DB) *gorm.DB {
	if len(values) == 0 {
		return nil
	}

	parts := strings.Split(values[0], ",")
	if len(parts) != 2 {
		return nil
	}

	startTime, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return nil
	}

	endTime, err := time.Parse("2006-01-02", parts[1])
	if err != nil {
		return nil
	}

	if filterType == FilterTypeDateTime {
		endTime = endTime.Add(24 * time.Hour).Add(-time.Second)
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), startTime, endTime)
	}
}

func (sb *ScopeBuilder) buildEnumScope(field string, values []string, allowedValues []string) func(*gorm.DB) *gorm.DB {
	if len(values) == 0 {
		return nil
	}

	// Validate enum value
	value := values[0]
	valid := false
	for _, allowed := range allowedValues {
		if value == allowed {
			valid = true
			break
		}
	}

	if !valid {
		return nil
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}
