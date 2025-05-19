package controller

import (
	"github.com/PhantomX7/dhamma/modules/event_attendance"
)

type controller struct {
	eventAttendanceService event_attendance.Service
}

func New(eventAttendanceService event_attendance.Service) event_attendance.Controller {
	return &controller{
		eventAttendanceService: eventAttendanceService,
	}
}
