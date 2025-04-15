package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/user_role"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.UserRole]
	db   *gorm.DB
}

func New(db *gorm.DB) user_role.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.UserRole](db),
		db:   db,
	}
}

func (r *repository) prepareDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}
