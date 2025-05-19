package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event_attendance"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.EventAttendance] // Use the interface type
	db   *gorm.DB
}

// New creates a new event_attendance repository instance.
func New(db *gorm.DB) event_attendance.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.EventAttendance](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all event_attendance entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.EventAttendance, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a event_attendance entity by its ID.
func (r *repository) FindByID(ctx context.Context, eventAttendanceID uint64, preloads ...string) (entity.EventAttendance, error) {
	return r.base.FindByID(ctx, eventAttendanceID, preloads...)
}

// Create creates a new event_attendance entity.
func (r *repository) Create(ctx context.Context, eventAttendance *entity.EventAttendance, tx *gorm.DB) error {
	return r.base.Create(ctx, eventAttendance, tx)
}

// Update updates an existing event_attendance entity.
func (r *repository) Update(ctx context.Context, eventAttendance *entity.EventAttendance, tx *gorm.DB) error {
	return r.base.Update(ctx, eventAttendance, tx)
}

// Delete deletes a event_attendance entity.
func (r *repository) Delete(ctx context.Context, eventAttendance *entity.EventAttendance, tx *gorm.DB) error {
	return r.base.Delete(ctx, eventAttendance, tx)
}

// Count counts event_attendance entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves event_attendance entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.EventAttendance, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single event_attendance entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.EventAttendance, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves event_attendance entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.EventAttendance, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single event_attendance entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.EventAttendance, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any event_attendance records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
