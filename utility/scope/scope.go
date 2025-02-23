package scope

import (
	"fmt"

	"gorm.io/gorm"
)

type Scope = func(*gorm.DB) *gorm.DB

// LimitScope will return a scope with limit
func LimitScope(limit int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// OffsetScope will return a scope with offset
func OffsetScope(limit int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(limit)
	}
}

// OrderScope will return a scope with order
func OrderScope(order string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

// WhereNotInScope will return a scope with WHERE NOT IN condition
// the argument must be slice of something
func WhereNotInScope(key string, value interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` NOT IN ?", key), value)
	}
}

// WhereNotInScope will return a scope with WHERE IN condition
// the argument must be slice of something
func WhereInScope(key string, value interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` IN ?", key), value)
	}
}

// WhereIsScope will return a scope with = condition
func WhereIsScope(key string, value interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` = ?", key), value)
	}
}

// WhereIsNotScope will return a scope with <> condition
func WhereIsNotScope(key string, value interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` <> ?", key), value)
	}
}

// WhereIsScope will return a scope with = condition
func WhereLikeScope(key string, value string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` LIKE ?", key), "%"+value+"%")
	}
}

// WhereBetweenScope will return a scope with BETWEEN value1 AND value2 condition
func WhereBetweenScope(key string, value1 interface{}, value2 interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` BETWEEN ? AND ?", key), value1, value2)
	}
}

// WhereIsNull will return a scope with null value for given key
func WhereIsNullScope(key string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` IS NULL", key))
	}
}

// WhereIsNotNull will return a scope with not null value for given key
func WhereIsNotNullScope(key string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` IS NOT NULL", key))
	}
}

// WhereIsNotNull will return a scope with not null value for given key
func JoinScope(key string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(key)
	}
}

// WhereIsNotNull will return a scope with not null value for given key
func PreloadScope(key string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(key)
	}
}
