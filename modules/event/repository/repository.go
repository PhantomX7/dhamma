package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.Event] // Use the interface type
	db   *gorm.DB
}

// New creates a new event repository instance.
func New(db *gorm.DB) event.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.Event](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all event entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Event, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a event entity by its ID.
func (r *repository) FindByID(ctx context.Context, eventID uint64, preloads ...string) (entity.Event, error) {
	return r.base.FindByID(ctx, eventID, preloads...)
}

// Create creates a new event entity.
func (r *repository) Create(ctx context.Context, event *entity.Event, tx *gorm.DB) error {
	return r.base.Create(ctx, event, tx)
}

// Update updates an existing event entity.
func (r *repository) Update(ctx context.Context, event *entity.Event, tx *gorm.DB) error {
	return r.base.Update(ctx, event, tx)
}

// Delete deletes a event entity.
func (r *repository) Delete(ctx context.Context, event *entity.Event, tx *gorm.DB) error {
	return r.base.Delete(ctx, event, tx)
}

// Count counts event entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves event entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.Event, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single event entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.Event, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves event entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.Event, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single event entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.Event, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any event records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
