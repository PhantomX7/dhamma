package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/utility/pagination"
	"gorm.io/gorm"
)

// BaseRepositoryInterface defines common repository operations
type BaseRepositoryInterface[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, model *T, tx *gorm.DB) error

	// Update updates an existing entity
	Update(ctx context.Context, model *T, tx *gorm.DB) error

	// Delete deletes an entity
	Delete(ctx context.Context, model *T, tx *gorm.DB) error

	// FindByID retrieves an entity by ID
	FindByID(ctx context.Context, id uint64, preloads ...string) (T, error)

	// FindAll retrieves all entities with pagination
	FindAll(ctx context.Context, pg *pagination.Pagination) ([]T, error)

	// Count counts entities with pagination filters
	Count(ctx context.Context, pg *pagination.Pagination) (int64, error)

	// FindByField retrieves entities where a specific field matches the given value.
	FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]T, error)

	// FindOneByField retrieves a single entity where a specific field matches the given value.
	FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (T, error)

	// FindByFields retrieves entities matching multiple field conditions.
	FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]T, error)

	// FindOneByFields retrieves a single entity matching multiple field conditions.
	FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (T, error)

	// Exists checks if any records match the given conditions.
	Exists(ctx context.Context, conditions map[string]any) (bool, error)
}
