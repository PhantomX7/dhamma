package service

import (
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
)

type service struct {
	userRepo         user.Repository
	refreshTokenRepo refresh_token.Repository
}

func New(
	userRepo user.Repository,
	refreshTokenRepo refresh_token.Repository,
) user.Service {
	return &service{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}
