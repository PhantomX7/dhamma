package middleware

import "github.com/PhantomX7/dhamma/modules/refresh_token"

type Middleware struct {
	refreshTokenRepo refresh_token.Repository
	// enforcer       *casbin.Enforcer
}

func New(refreshTokenRepo refresh_token.Repository) *Middleware {
	return &Middleware{
		refreshTokenRepo: refreshTokenRepo,
	}
}
