package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.Domain]
	db   *gorm.DB
}

func New(db *gorm.DB) domain.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Domain](db),
		db:   db,
	}
}

func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Domain, error) {
	return r.base.FindAll(ctx, pg)
}

func (r *repository) FindByID(ctx context.Context, domainID uint64, preloads ...string) (entity.Domain, error) {
	return r.base.FindByID(ctx, domainID, preloads...)
}

func (r *repository) Create(ctx context.Context, domain *entity.Domain, tx *gorm.DB) error {
	return r.base.Create(ctx, domain, tx)
}

func (r *repository) Update(ctx context.Context, domain *entity.Domain, tx *gorm.DB) error {
	return r.base.Update(ctx, domain, tx)
}

func (r *repository) Delete(ctx context.Context, domain *entity.Domain, tx *gorm.DB) error {
	return r.base.Delete(ctx, domain, tx)
}

func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}
