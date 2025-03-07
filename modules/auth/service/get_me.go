package service

import (
	"context"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/dhamma/utility"
)

func (u *service) GetMe(ctx context.Context) (res response.MeResponse, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	user, err := u.userRepo.FindByID(ctx, contextValues.UserID, true)
	if err != nil {
		return
	}

	// override user roles with specific domain roles
	if contextValues.DomainID != nil {
		userRoles, err := u.userRoleRepo.FindByUserIDAndDomainID(ctx, contextValues.UserID, *contextValues.DomainID, true)
		if err != nil {
			return response.MeResponse{}, err
		}
		user.UserRoles = userRoles
	}

	res.User = user

	return
}
