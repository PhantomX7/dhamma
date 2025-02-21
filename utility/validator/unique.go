package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// check if value of request is unique in database
// tag format : unique=tablename.columnname
func (cv cValidator) Unique() validator.Func {
	return func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Param(), ".")
		var count int64
		query := cv.db.Table(arr[0]).Where("`"+arr[1]+"` = ?", fl.Field().Interface())
		if cv.db.Migrator().HasColumn(arr[0], "deleted_at") {
			query = query.Where("`deleted_at` IS NULL")
		}
		err := query.Count(&count).Error
		return err == nil && count == 0
	}
}
