package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
)

func (u *service) UpdatePassword(
	userID uint64,
	request request.UpdatePasswordRequest,
	ctx context.Context,
) (err error) {
	userM, err := u.userRepo.FindByID(userID, false, ctx)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userM.Password), []byte(request.CurrentPassword))
	if err != nil {
		err = errors.New("current password is incorrect")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("failed to update password")
		return
	}
	userM.Password = string(password)

	err = u.userRepo.Update(&userM, nil, ctx)
	if err != nil {
		return
	}

	return
}
