package service

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) Create(ctx context.Context, request request.EventCreateRequest) (event entity.Event, err error) {
	_, err = utility.CheckDomainContext(ctx, request.DomainID, "event", "create")
	if err != nil {
		return
	}

	err = copier.Copy(&event, &request)
	if err != nil {
		return
	}

	err = s.eventRepo.Create(ctx, &event, nil)
	if err != nil {
		return
	}

	return
}
