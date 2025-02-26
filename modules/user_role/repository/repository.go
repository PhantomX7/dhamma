package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/user_role"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) user_role.Repository {
	return &repository{
		db: db,
	}
}
