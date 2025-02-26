package service

import (
	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const AccessTokenExpiry = 30 * time.Minute

func (u *service) GenerateAccessToken(userID uint64, role string) (string, error) {
	claims := entity.AccessClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
