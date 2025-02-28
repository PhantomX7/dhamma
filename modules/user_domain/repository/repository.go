package repository

import (
	"github.com/PhantomX7/dhamma/modules/user_domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) user_domain.Repository {
	return &repository{
		db: db,
	}
}
