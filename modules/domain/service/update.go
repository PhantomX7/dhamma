package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/domain/dto/request"
)

func (s *service) Update(domainID uint64, request request.DomainUpdateRequest, ctx context.Context) (domain entity.Domain, err error) {
	domain, err = s.domainRepo.FindByID(domainID, ctx)
	if err != nil {
		return
	}

	err = copier.Copy(&domain, &request)
	if err != nil {
		return
	}

	err = s.domainRepo.Update(&domain, nil, ctx)
	if err != nil {
		return
	}

	return
}
