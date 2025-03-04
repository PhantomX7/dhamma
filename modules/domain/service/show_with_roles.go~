package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (s *service) ShowWithRoles(domainID uint64, ctx context.Context) (domain entity.Domain, err error) {
	domain, err = s.domainRepo.GetDomainRoles(domainID, ctx)
	if err != nil {
		return
	}

	return
}
