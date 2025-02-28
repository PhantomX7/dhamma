package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements domain.Service
func (s *service) Show(domainID uint64, ctx context.Context) (domain entity.Domain, err error) {
	domain, err = s.domainRepo.FindByID(domainID, ctx)
	if err != nil {
		return
	}

	return
}
