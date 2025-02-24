package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/refresh_token"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) refresh_token.Repository {
	return &repository{
		db: db,
	}
}
