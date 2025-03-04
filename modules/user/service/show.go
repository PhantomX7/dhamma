package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements user.Service.
func (s *service) Show(userID uint64, ctx context.Context) (user entity.User, err error) {
	haveDomain, domainID := utility.GetDomainIDFromContext(ctx)

	user, err = s.userRepo.FindByID(userID, true, ctx)
	if err != nil {
		return
	}

	if haveDomain {
		var validDomain bool
		validDomain, err = s.userDomainRepo.HasDomain(userID, domainID, ctx)
		if err != nil {
			return
		}

		if !validDomain {
			err = errors.New("forbidden")
			return
		}

		var userRoles []entity.UserRole
		userRoles, err = s.userRoleRepo.FindByUserIDAndDomainID(userID, domainID, true, ctx)
		if err != nil {
			return
		}
		user.UserRoles = userRoles

		var domain entity.Domain
		domain, err = s.domainRepo.FindByID(domainID, ctx)
		if err != nil {
			return
		}
		user.Domains = []entity.Domain{domain}
	}

	return
}
