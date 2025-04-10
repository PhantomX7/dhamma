package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.Role]
	db   *gorm.DB
}

func New(db *gorm.DB) role.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Role](db),
		db:   db,
	}
}

func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Role, error) {
	return r.base.FindAll(ctx, pg)
}

func (r *repository) FindByID(ctx context.Context, roleID uint64, preloads ...string) (entity.Role, error) {
	return r.base.FindByID(ctx, roleID, preloads...)
}

func (r *repository) Create(ctx context.Context, role *entity.Role, tx *gorm.DB) error {
	return r.base.Create(ctx, role, tx)
}

func (r *repository) Update(ctx context.Context, role *entity.Role, tx *gorm.DB) error {
	return r.base.Update(ctx, role, tx)
}

func (r *repository) Delete(ctx context.Context, role *entity.Role, tx *gorm.DB) error {
	return r.base.Delete(ctx, role, tx)
}

func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}
