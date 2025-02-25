package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/golang-jwt/jwt/v4"
)

func (u *service) Refresh(request request.RefreshRequest, ctx context.Context) (res response.AuthResponse, err error) {
	// Parse refresh token and validate
	claims := &middleware.RefreshClaims{}
	token, err := jwt.ParseWithClaims(request.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.JWT_SECRET), nil
	})
	if err != nil {
		err = errors.New("invalid or expired token")
		return
	}

	if !token.Valid {
		err = errors.New("invalid token")
		return
	}

	// Check if refresh token is valid
	refreshTokenM, err := u.refreshTokenRepo.FindByID(claims.RefreshToken, ctx)
	if err != nil {
		err = errors.New("invalid refresh token")
		return
	}

	accessToken, err := middleware.GenerateAccessToken(refreshTokenM.UserID, "admin")
	if err != nil {
		return
	}

	refreshTokenM.IsValid = false

	tx := u.transactionManager.NewTransaction()

	err = u.refreshTokenRepo.Update(&refreshTokenM, tx, ctx)
	if err != nil {
		tx.Rollback()
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(refreshTokenM.UserID, tx, u.refreshTokenRepo)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	res = response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}
