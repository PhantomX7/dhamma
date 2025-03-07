package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/permission"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) permission.Repository {
	return &repository{
		db: db,
	}
}
