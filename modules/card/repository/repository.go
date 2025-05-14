package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/card"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.Card] // Use the interface type
	db   *gorm.DB
}

// New creates a new card repository instance.
func New(db *gorm.DB) card.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Card](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all card entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Card, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a card entity by its ID.
func (r *repository) FindByID(ctx context.Context, cardID uint64, preloads ...string) (entity.Card, error) {
	return r.base.FindByID(ctx, cardID, preloads...)
}

// Create creates a new card entity.
func (r *repository) Create(ctx context.Context, card *entity.Card, tx *gorm.DB) error {
	return r.base.Create(ctx, card, tx)
}

// Update updates an existing card entity.
func (r *repository) Update(ctx context.Context, card *entity.Card, tx *gorm.DB) error {
	return r.base.Update(ctx, card, tx)
}

// Delete deletes a card entity.
func (r *repository) Delete(ctx context.Context, card *entity.Card, tx *gorm.DB) error {
	return r.base.Delete(ctx, card, tx)
}

// Count counts card entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves card entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.Card, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single card entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.Card, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves card entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.Card, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single card entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.Card, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any card records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
