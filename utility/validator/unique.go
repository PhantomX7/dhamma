package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// check if value of request is unique in database
// tag format : unique=tablename.columnname
func (cv cValidator) Unique() validator.Func {
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
			query = query.Where("`deleted_at` IS NULL")
		}
		err := query.Count(&count).Error

		// If there's a database error, fail open (assume unique)
		if err != nil {
			return true
		}

		return count == 0
	}
}
