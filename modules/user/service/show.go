package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements user.Service.
func (s *service) Show(ctx context.Context, userID uint64) (user entity.User, err error) {
	haveDomain, domainID := utility.GetDomainIDFromContext(ctx)

	user, err = s.userRepo.FindByID(ctx, userID, true)
	if err != nil {
		return
	}

	if haveDomain {
		var validDomain bool
		validDomain, err = s.userDomainRepo.HasDomain(ctx, userID, domainID)
		if err != nil {
			return
		}

		if !validDomain {
			err = errors.New("forbidden")
			return
		}

		var userRoles []entity.UserRole
		userRoles, err = s.userRoleRepo.FindByUserIDAndDomainID(ctx, userID, domainID, true)
		if err != nil {
			return
		}
		user.UserRoles = userRoles

		var domain entity.Domain
		domain, err = s.domainRepo.FindByID(ctx, domainID)
		if err != nil {
			return
		}
		user.Domains = []entity.Domain{domain}
	}

	return
}
