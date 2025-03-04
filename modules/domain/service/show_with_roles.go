package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (s *service) ShowWithRoles(ctx context.Context, domainID uint64) (domain entity.Domain, err error) {
	domain, err = s.domainRepo.GetDomainRoles(ctx, domainID)
	if err != nil {
		return
	}

	return
}
