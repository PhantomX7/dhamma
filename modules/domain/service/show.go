package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements domain.Service
func (s *service) Show(ctx context.Context, domainID uint64) (domain entity.Domain, err error) {
	domain, err = s.domainRepo.FindByID(ctx, domainID)
	if err != nil {
		return
	}

	return
}
