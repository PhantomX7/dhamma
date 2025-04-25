package service

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

// Show implements user.Service.
func (s *service) Show(ctx context.Context, userID uint64) (user entity.User, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	user, err = s.userRepo.FindByID(ctx, userID, "Domains", "UserRoles.Role")
	if err != nil {
		return
	}

	if contextValues.DomainID != nil {
		var validDomain bool
		validDomain, err = s.userDomainRepo.HasDomain(ctx, userID, *contextValues.DomainID)
		if err != nil {
			return
		}

		if !validDomain {
			err = errors.New("forbidden")
			return
		}

		var userRoles []entity.UserRole
		userRoles, err = s.userRoleRepo.FindByUserIDAndDomainID(
			ctx,
			userID,
			*contextValues.DomainID,
			"Domain", "Role",
		)
		if err != nil {
			return
		}

		for _, userRole := range userRoles {
			userRole.Role.Permissions = s.casbin.GetRolePermissions(userRole.Role.ID, userRole.Role.DomainID)
		}
		user.UserRoles = userRoles

		var domain entity.Domain
		domain, err = s.domainRepo.FindByID(ctx, *contextValues.DomainID)
		if err != nil {
			return
		}
		user.Domains = []entity.Domain{domain}
	}

	return
}
