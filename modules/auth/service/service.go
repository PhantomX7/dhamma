package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/PhantomX7/dhamma/config"
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
	ID uint64 `json:"id"`
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

func generateTokenByID(id uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: jwt.TimeFunc().Add(time.Hour * 24),
			},
			IssuedAt: &jwt.NumericDate{
				Time: jwt.TimeFunc(),
			},
		},
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		err = errors.New("failed to generate token")
		return "", err
	}

	return tokenString, nil
}
