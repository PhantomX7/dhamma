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

// Validator interface for direct method calls in tests
type Validator interface {
	CheckUnique(value, table, column string) bool
	CheckExist(value, table, column string) bool
}

type cValidator struct {
	validator *validator.Validate
	db        *gorm.DB
}

func New(db *gorm.DB) CustomValidator {
	return cValidator{
		validator: validator.New(),
		db:        db,
	}
}

// NewValidator creates a new validator that implements the Validator interface
func NewValidator(db *gorm.DB) Validator {
	return &cValidator{
		validator: validator.New(),
		db:        db,
	}
}

// UniqueValue checks if a value is unique in the specified table and column
func (cv *cValidator) UniqueValue(value, table, column string) bool {
	var count int64

	query := cv.db.Table(table).Where("`"+column+"` = ?", value)
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

// ExistValue checks if a value exists in the specified table and column
func (cv *cValidator) ExistValue(value, table, column string) bool {
	var count int64

	query := cv.db.Table(table).Where("`"+column+"` = ?", value)
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

// CheckUnique method for Validator interface - delegates to UniqueValue
func (cv *cValidator) CheckUnique(value, table, column string) bool {
	return cv.UniqueValue(value, table, column)
}

// CheckExist method for Validator interface - delegates to ExistValue
func (cv *cValidator) CheckExist(value, table, column string) bool {
	return cv.ExistValue(value, table, column)
}
