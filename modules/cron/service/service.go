package service

import (
	"github.com/PhantomX7/dhamma/modules/cron"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
)

type service struct {
	refreshTokenRepo refresh_token.Repository
}

func New(
	refreshTokenRepo refresh_token.Repository,
) cron.Service {
	return &service{
		refreshTokenRepo: refreshTokenRepo,
	}
}
