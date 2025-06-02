package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// check if value of request exist in database
// tag format : exist=tablename.columnname
func (cv cValidator) Exist() validator.Func {
	return func(fl validator.FieldLevel) bool {
		var count int64

		arr := strings.Split(fl.Param(), ".")
		// Validate parameter format
		if len(arr) != 2 {
			return true // Invalid format, validation passes (fail open)
		}

		table, column := arr[0], arr[1]
		query := cv.db.Table(table).Where("`"+column+"` = ?", fl.Field().Interface())
		if cv.db.Migrator().HasColumn(table, "deleted_at") {
			query = query.Not("deleted_at IS NOT NULL")
		}
		err := query.Count(&count).Error

		// If there's a database error, fail closed (assume doesn't exist)
		if err != nil {
			return false
		}

		return count > 0
	}
}
