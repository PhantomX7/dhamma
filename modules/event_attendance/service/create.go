package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event_attendance/dto/request"
)

func (s *service) Create(ctx context.Context, request request.EventAttendanceCreateRequest) (eventAttendance entity.EventAttendance, err error) {
	eventAttendance = entity.EventAttendance{
		IsActive: true,
	}

	err = copier.Copy(&eventAttendance, &request)
	if err != nil {
		return
	}

	err = s.eventAttendanceRepo.Create(ctx, &eventAttendance, nil)
	if err != nil {
		return
	}

	return
}
