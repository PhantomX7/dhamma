package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event_attendance/dto/request"
)

func (s *service) Update(ctx context.Context, eventAttendanceID uint64, request request.EventAttendanceUpdateRequest) (eventAttendance entity.EventAttendance, err error) {
	eventAttendance, err = s.eventAttendanceRepo.FindByID(ctx, eventAttendanceID)
	if err != nil {
		return
	}

	err = copier.Copy(&eventAttendance, &request)
	if err != nil {
		return
	}

	err = s.eventAttendanceRepo.Update(ctx, &eventAttendance, nil)
	if err != nil {
		return
	}

	return
}
