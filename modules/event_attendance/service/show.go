package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements event_attendance.Service
func (s *service) Show(ctx context.Context, eventAttendanceID uint64) (eventAttendance entity.EventAttendance, err error) {
	eventAttendance, err = s.eventAttendanceRepo.FindByID(ctx, eventAttendanceID)
	if err != nil {
		return
	}

	return
}
