package service

import (
	"context"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/dhamma/utility"
)

func (u *service) GetMe(userID uint64, ctx context.Context) (res response.MeResponse, err error) {
	hasDomain, domainID := utility.GetDomainIDFromContext(ctx)

	user, err := u.userRepo.FindByID(userID, true, ctx)
	if err != nil {
		return
	}

	if hasDomain {
		userRoles, err := u.userRoleRepo.FindByUserIDAndDomainID(userID, domainID, true, ctx)
		if err != nil {
			return response.MeResponse{}, err
		}
		user.UserRoles = userRoles
	}

	res.User = user

	return
}
