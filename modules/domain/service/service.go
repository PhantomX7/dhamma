package service

import (
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
)

type service struct {
	domainRepo       domain.Repository
	refreshTokenRepo refresh_token.Repository
}

func New(
	domainRepo domain.Repository,
	refreshTokenRepo refresh_token.Repository,
) domain.Service {
	return &service{
		domainRepo:       domainRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}
