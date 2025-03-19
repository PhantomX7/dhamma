package middleware

import (
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_domain"

	"go.uber.org/zap"
)

type Middleware struct {
	userRepo         user.Repository
	refreshTokenRepo refresh_token.Repository
	userDomainRepo   user_domain.Repository
	domainRepo       domain.Repository
	casbin           casbin.Client
	logger           *zap.Logger
}

func New(
	userRepo user.Repository,
	refreshTokenRepo refresh_token.Repository,
	userDomainRepo user_domain.Repository,
	domainRepo domain.Repository,
	casbin casbin.Client,
	logger *zap.Logger,
) *Middleware {
	return &Middleware{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		userDomainRepo:   userDomainRepo,
		domainRepo:       domainRepo,
		casbin:           casbin,
		logger:           logger,
	}
}
