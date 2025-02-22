package service

import (
	"github.com/PhantomX7/dhamma/modules/user"
)

type service struct {
	userRepo user.Repository
}

func New(
	userRepo user.Repository,
) user.Service {
	return &service{
		userRepo: userRepo,
	}
}
