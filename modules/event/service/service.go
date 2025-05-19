package service

import (
	"github.com/PhantomX7/dhamma/modules/event"
)

type service struct {
	eventRepo event.Repository
}

func New(
	eventRepo event.Repository,
) event.Service {
	return &service{
		eventRepo: eventRepo,
	}
}
