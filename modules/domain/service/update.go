package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
)

func (s *service) Update(ctx context.Context, domainID uint64, request request.DomainUpdateRequest) (domain entity.Domain, err error) {
	domain, err = s.domainRepo.FindByID(ctx, domainID)
	if err != nil {
		return
	}

	err = copier.Copy(&domain, &request)
	if err != nil {
		return
	}

	err = s.domainRepo.Update(ctx, &domain, nil)
	if err != nil {
		return
	}

	return
}
