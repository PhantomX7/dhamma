package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

// Show implements event.Service
func (s *service) Show(ctx context.Context, eventID uint64) (event entity.Event, err error) {
	event, err = s.eventRepo.FindByID(ctx, eventID, "Domain")
	if err != nil {
		return
	}

	// Perform domain context check using the generic helper
	_, err = utility.CheckDomainContext(ctx, event.DomainID, "event", "show")
	if err != nil {
		return
	}

	return
}
