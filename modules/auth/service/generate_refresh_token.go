package service

import (
	"context"
	"time"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/entity"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const RefreshTokenExpiry = 3 * 24 * time.Hour

func (u *service) GenerateRefreshToken(userID uint64, tx *gorm.DB) (string, error) {
	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(RefreshTokenExpiry),
		IsValid:   true,
	}

	// Save to database
	err := u.refreshTokenRepo.Create(context.Background(), refreshToken, tx)
	if err != nil {
		return "", err
	}

	claims := entity.RefreshClaims{
		RefreshToken: refreshToken.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
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
