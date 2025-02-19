package validators

import (
	"database/sql"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var CustomValidator *cValidator

type cValidator struct {
	db *gorm.DB
}

func NewValidator(db *gorm.DB) {
	CustomValidator = &cValidator{
		db: db,
	}
}

// check if value of request is unique in database
// tag format : unique=tablename.columnname
func (cv *cValidator) Unique() validator.Func {
	return func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Param(), ".")
		var deletedAt sql.NullString
		err := cv.db.Table(arr[0]).Select("deleted_at").Where("`"+arr[1]+"` = ?", fl.Field().Interface()).First(&deletedAt).Error
		if err != nil {
			return false
		}
		return deletedAt.Valid
	}
}

// check if value of request exist in database
// tag format : exist=tablename.columnname
func (cv *cValidator) Exist() validator.Func {
	return func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Param(), ".")
		if len(arr) != 2 {
			return false
		}

		var count int64
		err := cv.db.Table(arr[0]).Where("`"+arr[1]+"` = ? AND `deleted_at` IS NOT NULL", fl.Field().Interface()).Count(&count).Error
		return err == nil && count > 0
	}
}
