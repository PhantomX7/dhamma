package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/user"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.Repository {
	return &repository{
		db: db,
	}
}
