package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.RefreshToken]
	db   *gorm.DB
}

func New(db *gorm.DB) refresh_token.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.RefreshToken](db),
		db:   db,
	}
}

func (r *repository) Create(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error {
	return r.base.Create(ctx, refreshToken, tx)
}

func (r *repository) Update(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error {
	return r.base.Update(ctx, refreshToken, tx)
}

func (r *repository) Delete(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error {
	return r.base.Delete(ctx, refreshToken, tx)
}
