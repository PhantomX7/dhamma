package repository

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/user_domain"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.UserDomain]
	db   *gorm.DB
}

func New(db *gorm.DB) user_domain.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.UserDomain](db),
		db:   db,
	}
}
