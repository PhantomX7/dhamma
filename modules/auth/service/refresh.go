package service

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (u *service) Refresh(request request.RefreshRequest, ctx context.Context) (res response.AuthResponse, err error) {
	refreshTokenM, err := u.refreshTokenRepo.FindByID(request.RefreshToken, ctx)
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
		RefreshToken: refreshToken.ID.String(),
	}
	return
}
