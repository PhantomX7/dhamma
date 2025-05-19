package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/point_mutation"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.PointMutation] // Use the interface type
	db   *gorm.DB
}

// New creates a new point_mutation repository instance.
func New(db *gorm.DB) point_mutation.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.PointMutation](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all point_mutation entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.PointMutation, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a point_mutation entity by its ID.
func (r *repository) FindByID(ctx context.Context, pointMutationID uint64, preloads ...string) (entity.PointMutation, error) {
	return r.base.FindByID(ctx, pointMutationID, preloads...)
}

// Create creates a new point_mutation entity.
func (r *repository) Create(ctx context.Context, pointMutation *entity.PointMutation, tx *gorm.DB) error {
	return r.base.Create(ctx, pointMutation, tx)
}

// Update updates an existing point_mutation entity.
func (r *repository) Update(ctx context.Context, pointMutation *entity.PointMutation, tx *gorm.DB) error {
	return r.base.Update(ctx, pointMutation, tx)
}

// Delete deletes a point_mutation entity.
func (r *repository) Delete(ctx context.Context, pointMutation *entity.PointMutation, tx *gorm.DB) error {
	return r.base.Delete(ctx, pointMutation, tx)
}

// Count counts point_mutation entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves point_mutation entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.PointMutation, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single point_mutation entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.PointMutation, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves point_mutation entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.PointMutation, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single point_mutation entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.PointMutation, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any point_mutation records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
