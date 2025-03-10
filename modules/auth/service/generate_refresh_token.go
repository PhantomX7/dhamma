package service

import (
	"context"
	"time"

	"github.com/PhantomX7/dhamma/constants"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/entity"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *service) GenerateRefreshToken(userID uint64, tx *gorm.DB) (string, error) {
	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(constants.RefreshTokenExpiry),
		IsValid:   true,
	}

	// Save to database
	err := s.refreshTokenRepo.Create(context.Background(), refreshToken, tx)
	if err != nil {
		return "", err
	}

	claims := entity.RefreshClaims{
		RefreshToken: refreshToken.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshToken.ExpiresAt),
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
