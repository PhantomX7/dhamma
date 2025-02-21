package validator

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// var CustomValidator *cValidator

type CustomValidator interface {
	Unique() validator.Func
	Exist() validator.Func
}

type cValidator struct {
	db *gorm.DB
}

func New(db *gorm.DB) CustomValidator {
	return cValidator{
		db: db,
	}
}
