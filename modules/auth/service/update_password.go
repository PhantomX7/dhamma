package service

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/utility"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
)

func (s *service) UpdatePassword(
	ctx context.Context,
	request request.UpdatePasswordRequest,
) (err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	userM, err := s.userRepo.FindByID(ctx, contextValues.UserID)
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

	err = s.userRepo.Update(ctx, &userM, nil)
	if err != nil {
		return
	}

	return
}
