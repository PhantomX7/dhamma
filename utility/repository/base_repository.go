package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/utility/pagination"
	"gorm.io/gorm"
)

// BaseRepository provides common repository functionality with generics.
// It implements standard CRUD operations and query methods for any entity type.
type BaseRepository[T any] struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new base repository with the given database connection.
// It returns a repository instance that can perform operations on entities of type T.
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return BaseRepository[T]{
		DB: db,
	}
}

// WithTx returns a new repository instance that uses the provided transaction.
// If tx is nil, the original repository's DB connection is used.
// This allows for method chaining in transaction blocks.
func (r BaseRepository[T]) WithTx(tx *gorm.DB) BaseRepository[T] {
	if tx != nil {
		r.DB = tx
	}
	return r
}

// Count returns the number of records matching the given pagination criteria.
// It applies any filters from the pagination object but ignores limit and offset.
// Returns the count and any error encountered during the database operation.
func (r BaseRepository[T]) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	var count int64
	scopeBuilder := pagination.NewScopeBuilder(pg)
	scopes, _ := scopeBuilder.Build()

	err := r.DB.WithContext(ctx).
		Scopes(scopes...).
		Model((*T)(nil)).
		Count(&count).Error
	if err != nil {
		return 0, WrapError(ErrDatabase, "failed to count records")
	}
	return count, nil
}

// FindByID retrieves a single entity by its primary key ID.
// It accepts a context for cancellation and timeout control, an ID to search for,
// and optional preload relationships.
// Returns the found entity and nil error on success, or a zero value entity and an error if:
// - The record is not found (ErrNotFound)
// - A database error occurs (ErrDatabase)
func (r BaseRepository[T]) FindByID(ctx context.Context, id uint64, preloads ...string) (T, error) {
	var entity T
	db := r.prepareDB(ctx, nil)
	db = r.applyPreloads(db, preloads...)

	result := db.First(&entity, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity, WrapError(ErrNotFound, "entity with ID %d not found", id)
	}
	if result.Error != nil {
		return entity, WrapError(ErrDatabase, "failed to find entity with ID %d", id)
	}

	return entity, nil
}

// Create persists a new entity to the database.
// It accepts a context for cancellation and timeout control, a pointer to the model to create,
// and an optional transaction for atomic operations.
// Returns nil on success or an error if the database operation fails.
func (r BaseRepository[T]) Create(ctx context.Context, model *T, tx *gorm.DB) error {
	db := r.prepareDB(ctx, tx)

	err := db.Create(model).Error
	if err != nil {
		return WrapError(ErrDatabase, "failed to create record")
	}
	return nil
}

// Update modifies an existing entity in the database.
// It accepts a context for cancellation and timeout control, a pointer to the model to update,
// and an optional transaction for atomic operations.
// The model must have its primary key set to identify the record to update.
// Returns nil on success or an error if the database operation fails.
func (r BaseRepository[T]) Update(ctx context.Context, model *T, tx *gorm.DB) error {
	db := r.prepareDB(ctx, tx)

	err := db.Save(model).Error
	if err != nil {
		return WrapError(ErrDatabase, "failed to update record")
	}
	return nil
}

// Delete removes an entity from the database.
// It accepts a context for cancellation and timeout control, a pointer to the model to delete,
// and an optional transaction for atomic operations.
// The model must have its primary key set to identify the record to delete.
// Returns nil on success or an error if the database operation fails.
func (r BaseRepository[T]) Delete(ctx context.Context, model *T, tx *gorm.DB) error {
	db := r.prepareDB(ctx, tx)

	err := db.Delete(model).Error
	if err != nil {
		return WrapError(ErrDatabase, "failed to delete record")
	}
	return nil
}

// FindAll retrieves multiple entities with pagination.
// It accepts a context for cancellation and timeout control and a pagination object
// that can include filters, sorting, limit, and offset.
// Returns a slice of entities and nil error on success, or an empty slice and an error
// if the database operation fails.
func (r BaseRepository[T]) FindAll(ctx context.Context, pg *pagination.Pagination) ([]T, error) {
	entities := make([]T, 0)
	scopeBuilder := pagination.NewScopeBuilder(pg)
	scopes, metaScopes := scopeBuilder.Build()

	err := r.DB.WithContext(ctx).
		Scopes(scopes...).
		Scopes(metaScopes...).
		Find(&entities).Error
	if err != nil {
		return nil, WrapError(ErrDatabase, "failed to find records")
	}

	return entities, nil
}

// FindByField retrieves entities where a specific field matches the given value.
// It accepts a context for cancellation and timeout control, the field name to filter on,
// the value to match, and optional preload relationships.
// Returns a slice of matching entities and nil error on success, or an empty slice and an error
// if the database operation fails.
func (r BaseRepository[T]) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]T, error) {
	var entities []T
	db := r.prepareDB(ctx, nil)
	db = r.applyPreloads(db, preloads...)

	err := db.Where(fieldName+" = ?", value).Find(&entities).Error
	if err != nil {
		return nil, WrapError(ErrDatabase, "failed to find records with %s=%v", fieldName, value)
	}
	return entities, nil
}

// FindOneByField retrieves a single entity where a specific field matches the given value.
// It accepts a context for cancellation and timeout control, the field name to filter on,
// the value to match, and optional preload relationships.
// Returns the found entity and nil error on success, or a zero value entity and an error if:
// - The record is not found (ErrNotFound)
// - A database error occurs (ErrDatabase)
func (r BaseRepository[T]) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (T, error) {
	var entity T
	db := r.prepareDB(ctx, nil)
	db = r.applyPreloads(db, preloads...)

	result := db.Where(fieldName+" = ?", value).First(&entity)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity, WrapError(ErrNotFound, "entity with %s=%v not found", fieldName, value)
	}
	if result.Error != nil {
		return entity, WrapError(ErrDatabase, "failed to find entity with %s=%v", fieldName, value)
	}

	return entity, nil
}

// FindByFields retrieves entities matching multiple field conditions.
// It accepts a context for cancellation and timeout control, a map of field names to values,
// and optional preload relationships.
// Returns a slice of matching entities and nil error on success, or an empty slice and an error
// if the database operation fails.
func (r BaseRepository[T]) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]T, error) {
	var entities []T
	db := r.prepareDB(ctx, nil)
	db = r.applyPreloads(db, preloads...)

	err := db.Where(conditions).Find(&entities).Error
	if err != nil {
		return nil, WrapError(ErrDatabase, "failed to find records with conditions")
	}
	return entities, nil
}

// Exists checks if any records match the given conditions.
// It accepts a context for cancellation and timeout control and a map of field names to values.
// Returns true if at least one matching record exists, false otherwise, and any error
// encountered during the database operation.
func (r BaseRepository[T]) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model((*T)(nil)).Where(conditions).Count(&count).Error
	if err != nil {
		return false, WrapError(ErrDatabase, "failed to check if record exists")
	}
	return count > 0, nil
}

// applyPreloads is a helper method to apply preload relationships to a query
func (r BaseRepository[T]) applyPreloads(db *gorm.DB, preloads ...string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// prepareDB is a helper method to prepare the database connection with context and transaction
func (r BaseRepository[T]) prepareDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx)
}
