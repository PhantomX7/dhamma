package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/follower"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.Follower] // Use the interface type
	db   *gorm.DB
}

// New creates a new follower repository instance.
func New(db *gorm.DB) follower.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Follower](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all follower entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Follower, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a follower entity by its ID.
func (r *repository) FindByID(ctx context.Context, followerID uint64, preloads ...string) (entity.Follower, error) {
	return r.base.FindByID(ctx, followerID, preloads...)
}

// Create creates a new follower entity.
func (r *repository) Create(ctx context.Context, follower *entity.Follower, tx *gorm.DB) error {
	return r.base.Create(ctx, follower, tx)
}

// Update updates an existing follower entity.
func (r *repository) Update(ctx context.Context, follower *entity.Follower, tx *gorm.DB) error {
	return r.base.Update(ctx, follower, tx)
}

// Delete deletes a follower entity.
func (r *repository) Delete(ctx context.Context, follower *entity.Follower, tx *gorm.DB) error {
	return r.base.Delete(ctx, follower, tx)
}

// Count counts follower entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves follower entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.Follower, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single follower entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.Follower, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves follower entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.Follower, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single follower entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.Follower, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any follower records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
