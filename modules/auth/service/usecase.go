package service

import (
	"github.com/golang-jwt/jwt/v4"

	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/PhantomX7/dhamma/modules/user"
)

type service struct {
	userRepo user.Repository
	// casbin   casbin.Client
	// cache    gocache.Client
}

type authClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	ID       uint64 `json:"id"`
	IssuedAt int64  `json:"orig_iat,omitempty"`
	Role     string `json:"role"`
}

func New(
	userRepo user.Repository,
	// casbin casbin.Client,
	// cache gocache.Client,
) auth.Service {
	return &service{
		userRepo: userRepo,
		// casbin:   casbin,
		// cache:    cache,
	}
}
