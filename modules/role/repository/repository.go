package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/role"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) role.Repository {
	return &repository{
		db: db,
	}
}
