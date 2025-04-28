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
	base baseRepo.BaseRepositoryInterface[entity.Role] // Use the interface type
	db   *gorm.DB
}

// New creates a new role repository instance.
func New(db *gorm.DB) role.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Role](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all role entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Role, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a role entity by its ID.
func (r *repository) FindByID(ctx context.Context, roleID uint64, preloads ...string) (entity.Role, error) {
	return r.base.FindByID(ctx, roleID, preloads...)
}

// Create creates a new role entity.
func (r *repository) Create(ctx context.Context, role *entity.Role, tx *gorm.DB) error {
	return r.base.Create(ctx, role, tx)
}

// Update updates an existing role entity.
func (r *repository) Update(ctx context.Context, role *entity.Role, tx *gorm.DB) error {
	return r.base.Update(ctx, role, tx)
}

// Delete deletes a role entity.
func (r *repository) Delete(ctx context.Context, role *entity.Role, tx *gorm.DB) error {
	return r.base.Delete(ctx, role, tx)
}

// Count counts role entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves role entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.Role, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single role entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.Role, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves role entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.Role, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single role entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.Role, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any role records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
