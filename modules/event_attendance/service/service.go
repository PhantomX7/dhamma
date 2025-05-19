package service

import (
	"github.com/PhantomX7/dhamma/modules/event_attendance"
)

type service struct {
	eventAttendanceRepo event_attendance.Repository
}

func New(
	eventAttendanceRepo event_attendance.Repository,
) event_attendance.Service {
	return &service{
		eventAttendanceRepo: eventAttendanceRepo,
	}
}
