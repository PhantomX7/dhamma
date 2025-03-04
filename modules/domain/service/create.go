package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
)

func (s *service) Create(ctx context.Context, request request.DomainCreateRequest) (domain entity.Domain, err error) {
	domain = entity.Domain{
		IsActive: true,
	}

	err = copier.Copy(&domain, &request)
	if err != nil {
		return
	}

	err = s.domainRepo.Create(ctx, &domain, nil)
	if err != nil {
		return
	}

	return
}
