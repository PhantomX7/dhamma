package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// check if value of request exist in database
// tag format : exist=tablename.columnname
func (cv cValidator) Exist() validator.Func {
	return func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Param(), ".")
		var count int64
		db := cv.db.Model(arr[0]).Where("`"+arr[1]+"` = ?", fl.Field().Interface())
		if cv.db.Migrator().HasColumn(arr[0], "deleted_at") {
			db = db.Not("deleted_at IS NOT NULL")
		}
		db.Count(&count)
		return count > 0
	}
}
