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
	base baseRepo.BaseRepositoryInterface[entity.Domain] // Use the interface type
	db   *gorm.DB
}

// New creates a new domain repository instance.
func New(db *gorm.DB) domain.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Domain](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all domain entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Domain, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a domain entity by its ID.
func (r *repository) FindByID(ctx context.Context, domainID uint64, preloads ...string) (entity.Domain, error) {
	return r.base.FindByID(ctx, domainID, preloads...)
}

// Create creates a new domain entity.
func (r *repository) Create(ctx context.Context, domain *entity.Domain, tx *gorm.DB) error {
	return r.base.Create(ctx, domain, tx)
}

// Update updates an existing domain entity.
func (r *repository) Update(ctx context.Context, domain *entity.Domain, tx *gorm.DB) error {
	return r.base.Update(ctx, domain, tx)
}

// Delete deletes a domain entity.
func (r *repository) Delete(ctx context.Context, domain *entity.Domain, tx *gorm.DB) error {
	return r.base.Delete(ctx, domain, tx)
}

// Count counts domain entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves domain entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.Domain, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single domain entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.Domain, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves domain entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.Domain, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single domain entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.Domain, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any domain records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
