package service

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"
)

func (s *service) AssignDomain(ctx context.Context, userID uint64, request request.AssignDomainRequest) (err error) {
	hasDomain, err := s.userDomainRepo.HasDomain(ctx, userID, request.DomainID)
	if err != nil || hasDomain {
		return
	}

	err = s.userDomainRepo.AssignDomain(ctx, userID, request.DomainID, nil)
	if err != nil {
		return
	}

	return
}
