package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/permission"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.Permission]
	db   *gorm.DB
}

func New(db *gorm.DB) permission.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Permission](db),
		db:   db,
	}
}

func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Permission, error) {
	return r.base.FindAll(ctx, pg)
}

func (r *repository) FindByID(ctx context.Context, permissionID uint64) (entity.Permission, error) {
	return r.base.FindByID(ctx, permissionID)
}

func (r *repository) Create(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	return r.base.Create(ctx, permission, tx)
}

func (r *repository) Update(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	return r.base.Update(ctx, permission, tx)
}

func (r *repository) Delete(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	return r.base.Delete(ctx, permission, tx)
}

func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}
