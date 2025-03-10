package service

import (
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_role"
)

type service struct {
	userRepo           user.Repository
	userRoleRepo       user_role.Repository
	refreshTokenRepo   refresh_token.Repository
	transactionManager transaction_manager.Client
	casbin             casbin.Client
}

func New(
	userRepo user.Repository,
	userRoleRepo user_role.Repository,
	refreshTokenRepo refresh_token.Repository,
	transactionManager transaction_manager.Client,
	casbin casbin.Client,
) auth.Service {
	return &service{
		userRepo:           userRepo,
		userRoleRepo:       userRoleRepo,
		refreshTokenRepo:   refreshTokenRepo,
		transactionManager: transactionManager,
		casbin:             casbin,
	}
}
