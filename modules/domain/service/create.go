package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
)

func (s *service) Create(request request.DomainCreateRequest, ctx context.Context) (domain entity.Domain, err error) {
	domain = entity.Domain{
		IsActive: true,
	}

	err = copier.Copy(&domain, &request)
	if err != nil {
		return
	}

	err = s.domainRepo.Create(&domain, nil, ctx)
	if err != nil {
		return
	}

	return
}
