package service

import (
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_domain"
	"github.com/PhantomX7/dhamma/modules/user_role"
)

type service struct {
	userRepo           user.Repository
	roleRepo           role.Repository
	userDomainRepo     user_domain.Repository
	userRoleRepo       user_role.Repository
	domainRepo         domain.Repository
	refreshTokenRepo   refresh_token.Repository
	transactionManager transaction_manager.Client
	casbin             casbin.Client
}

func New(
	userRepo user.Repository,
	roleRepo role.Repository,
	userDomainRepo user_domain.Repository,
	userRoleRepo user_role.Repository,
	domainRepo domain.Repository,
	refreshTokenRepo refresh_token.Repository,
	transactionManager transaction_manager.Client,
	casbin casbin.Client,
) user.Service {
	return &service{
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		userDomainRepo:     userDomainRepo,
		userRoleRepo:       userRoleRepo,
		domainRepo:         domainRepo,
		refreshTokenRepo:   refreshTokenRepo,
		transactionManager: transactionManager,
		casbin:             casbin,
	}
}
