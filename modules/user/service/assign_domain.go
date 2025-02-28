package service

import (
	"context"
	"github.com/PhantomX7/dhamma/modules/user/dto/request"
)

func (s *service) AssignDomain(userID uint64, request request.AssignDomainRequest, ctx context.Context) (err error) {
	hasDomain, err := s.userDomainRepo.HasDomain(userID, request.DomainID, ctx)
	if err != nil {
		return
	}

	if hasDomain {
		return
	}

	err = s.userDomainRepo.AssignDomain(userID, request.DomainID, nil, ctx)
	if err != nil {
		return
	}

	return
}
