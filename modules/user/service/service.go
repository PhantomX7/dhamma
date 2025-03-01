package service

import (
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_domain"
)

type service struct {
	userRepo           user.Repository
	userDomainRepo     user_domain.Repository
	refreshTokenRepo   refresh_token.Repository
	transactionManager transaction_manager.Client
}

func New(
	userRepo user.Repository,
	userDomainRepo user_domain.Repository,
	refreshTokenRepo refresh_token.Repository,
	transactionManager transaction_manager.Client,
) user.Service {
	return &service{
		userRepo:           userRepo,
		userDomainRepo:     userDomainRepo,
		refreshTokenRepo:   refreshTokenRepo,
		transactionManager: transactionManager,
	}
}
