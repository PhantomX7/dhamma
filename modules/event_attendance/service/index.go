package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements event_attendance.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	eventAttendances []entity.EventAttendance, meta utility.PaginationMeta, err error,
) {
	// Combine all scopes into a single AddCustomScope call
	pg.AddCustomScope(
		// Base join and preload
		func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("Follower").Preload("Event")
		},
	)

	eventAttendances, err = s.eventAttendanceRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.eventAttendanceRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
