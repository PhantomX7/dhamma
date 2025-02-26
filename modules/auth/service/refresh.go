package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/golang-jwt/jwt/v4"
)

func (u *service) Refresh(request request.RefreshRequest, ctx context.Context) (res response.AuthResponse, err error) {
	// Parse refresh token and validate
	claims := &entity.RefreshClaims{}
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

	accessToken, err := u.GenerateAccessToken(refreshTokenM.UserID, "admin")
	if err != nil {
		return
	}

	//invalidate current refresh token
	refreshTokenM.IsValid = false

	tx := u.transactionManager.NewTransaction()

	err = u.refreshTokenRepo.Update(&refreshTokenM, tx, ctx)
	if err != nil {
		tx.Rollback()
		return
	}

	refreshToken, err := u.GenerateRefreshToken(refreshTokenM.UserID, tx)
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
