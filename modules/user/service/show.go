package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements user.Service.
func (s *service) Show(userID uint64, ctx context.Context) (user entity.User, err error) {
	user, err = s.userRepo.FindByID(userID, true, ctx)
	if err != nil {
		return
	}

	return
}
