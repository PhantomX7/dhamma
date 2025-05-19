package service

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event/dto/request"
)

func (s *service) Update(ctx context.Context, eventID uint64, request request.EventUpdateRequest) (event entity.Event, err error) {
	event, err = s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return
	}

	err = copier.Copy(&event, &request)
	if err != nil {
		return
	}

	err = s.eventRepo.Update(ctx, &event, nil)
	if err != nil {
		return
	}

	return
}
