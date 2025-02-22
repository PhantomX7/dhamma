package service

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (u *service) GetMe(userID uint64, ctx context.Context) (res response.MeResponse, err error) {
	userM, err := u.userRepo.FindByID(userID, ctx)
	if err != nil {
		return
	}

	res.User = userM

	return
}
