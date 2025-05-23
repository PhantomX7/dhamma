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
		query := cv.db.Table(arr[0]).Where("`"+arr[1]+"` = ?", fl.Field().Interface())
		if cv.db.Migrator().HasColumn(arr[0], "deleted_at") {
			query = query.Not("deleted_at IS NOT NULL")
		}
		err := query.Count(&count).Error

		return err == nil && count > 0
	}
}
