package service

import (
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_domain"
)

type service struct {
	userRepo         user.Repository
	userDomainRepo   user_domain.Repository
	refreshTokenRepo refresh_token.Repository
}

func New(
	userRepo user.Repository,
	userDomainRepo user_domain.Repository,
	refreshTokenRepo refresh_token.Repository,
) user.Service {
	return &service{
		userRepo:         userRepo,
		userDomainRepo:   userDomainRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}
