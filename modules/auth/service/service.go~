package service

import (
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
)

type service struct {
	userRepo           user.Repository
	refreshTokenRepo   refresh_token.Repository
	transactionManager transaction_manager.Client
}

func New(
	userRepo user.Repository,
	refreshTokenRepo refresh_token.Repository,
	transactionManager transaction_manager.Client,
) auth.Service {
	return &service{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		transactionManager: transactionManager,
	}
}
