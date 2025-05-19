package controller

import (
	"github.com/PhantomX7/dhamma/modules/event"
)

type controller struct {
	eventService event.Service
}

func New(eventService event.Service) event.Controller {
	return &controller{
		eventService: eventService,
	}
}
