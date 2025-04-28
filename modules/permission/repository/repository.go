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
	base baseRepo.BaseRepositoryInterface[entity.Permission]
	db   *gorm.DB
}

// New creates a new permission repository instance.
func New(db *gorm.DB) permission.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Permission](db),
		db:   db,
	}
}

// FindAll retrieves all permission entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Permission, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a permission entity by its ID.
func (r *repository) FindByID(ctx context.Context, permissionID uint64, preloads ...string) (entity.Permission, error) {
	return r.base.FindByID(ctx, permissionID, preloads...)
}

// Create creates a new permission entity.
func (r *repository) Create(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	return r.base.Create(ctx, permission, tx)
}

// Update updates an existing permission entity.
func (r *repository) Update(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	return r.base.Update(ctx, permission, tx)
}

// Delete deletes a permission entity.
func (r *repository) Delete(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	return r.base.Delete(ctx, permission, tx)
}

// Count counts permission entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves permission entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.Permission, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single permission entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.Permission, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves permission entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.Permission, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single permission entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.Permission, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any permission records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
