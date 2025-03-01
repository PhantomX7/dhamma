package middleware

import (
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user_domain"
)

type Middleware struct {
	refreshTokenRepo refresh_token.Repository
	userDomainRepo   user_domain.Repository
	domainRepo       domain.Repository
}

func New(
	refreshTokenRepo refresh_token.Repository,
	userDomainRepo user_domain.Repository,
	domainRepo domain.Repository,
) *Middleware {
	return &Middleware{
		refreshTokenRepo: refreshTokenRepo,
		userDomainRepo:   userDomainRepo,
		domainRepo:       domainRepo,
	}
}
