package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements event.Service
func (s *service) Show(ctx context.Context, eventID uint64) (event entity.Event, err error) {
	event, err = s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return
	}

	return
}
