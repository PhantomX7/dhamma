package service

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) GetMe(ctx context.Context) (res response.MeResponse, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	user, err := s.userRepo.FindByID(ctx, contextValues.UserID, "Domains", "UserRoles.Role")
	if err != nil {
		return
	}

	// override user roles with specific domain roles
	if contextValues.DomainID != nil {
		userRoles, err := s.userRoleRepo.FindByUserIDAndDomainID(ctx, contextValues.UserID, *contextValues.DomainID, "Domain", "Role")
		if err != nil {
			return response.MeResponse{}, err
		}
		user.UserRoles = userRoles

		user.Permissions = s.casbin.GetUserPermissions(user.ID, *contextValues.DomainID)
	}

	res.User = user

	return
}
