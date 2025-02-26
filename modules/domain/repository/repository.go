package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/domain"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.Repository {
	return &repository{
		db: db,
	}
}
