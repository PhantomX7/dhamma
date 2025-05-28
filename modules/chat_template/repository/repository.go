package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/chat_template"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.ChatTemplate] // Use the interface type
	db   *gorm.DB
}

// New creates a new chat template repository instance.
func New(db *gorm.DB) chat_template.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.ChatTemplate](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all chat template entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.ChatTemplate, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a chat template entity by its ID.
func (r *repository) FindByID(ctx context.Context, chatTemplateID uint64, preloads ...string) (entity.ChatTemplate, error) {
	return r.base.FindByID(ctx, chatTemplateID, preloads...)
}

// Create creates a new chat template entity.
func (r *repository) Create(ctx context.Context, chatTemplate *entity.ChatTemplate, tx *gorm.DB) error {
	return r.base.Create(ctx, chatTemplate, tx)
}

// Update updates an existing chat template entity.
func (r *repository) Update(ctx context.Context, chatTemplate *entity.ChatTemplate, tx *gorm.DB) error {
	return r.base.Update(ctx, chatTemplate, tx)
}

// Delete deletes a chat template entity.
func (r *repository) Delete(ctx context.Context, chatTemplate *entity.ChatTemplate, tx *gorm.DB) error {
	return r.base.Delete(ctx, chatTemplate, tx)
}

// Count counts chat template entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves chat template entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.ChatTemplate, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single chat template entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.ChatTemplate, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves chat template entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.ChatTemplate, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single chat template entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.ChatTemplate, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any chat template records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
