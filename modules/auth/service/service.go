package service

import (
	"github.com/PhantomX7/dhamma/libs/database_transaction"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
)

type service struct {
	userRepo           user.Repository
	refreshTokenRepo   refresh_token.Repository
	transactionManager database_transaction.TransactionManager
	// casbin   casbin.Client
	// cache    gocache.Client
}

func New(
	userRepo user.Repository,
	refreshTokenRepo refresh_token.Repository,
	transactionManager database_transaction.TransactionManager,
	// casbin casbin.Client,
	// cache gocache.Client,
) auth.Service {
	return &service{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		transactionManager: transactionManager,
		// casbin:   casbin,
		// cache:    cache,
	}
}
