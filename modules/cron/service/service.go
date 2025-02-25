package service

import (
	"github.com/PhantomX7/dhamma/libs/database_transaction"
	"github.com/PhantomX7/dhamma/modules/cron"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
)

type service struct {
	refreshTokenRepo   refresh_token.Repository
	transactionManager database_transaction.TransactionManager
	// casbin   casbin.Client
	// cache    gocache.Client
}

func New(
	refreshTokenRepo refresh_token.Repository,
	transactionManager database_transaction.TransactionManager,
	// casbin casbin.Client,
	// cache gocache.Client,
) cron.Service {
	return &service{
		refreshTokenRepo:   refreshTokenRepo,
		transactionManager: transactionManager,
		// casbin:   casbin,
		// cache:    cache,
	}
}
