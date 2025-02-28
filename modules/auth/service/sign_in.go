package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (u *service) SignIn(request request.SignInRequest, ctx context.Context) (res response.AuthResponse, err error) {
	userM := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	userM, err = u.userRepo.FindByUsername(request.Username, ctx)
	if err != nil {
		err = errors.New("invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userM.Password), []byte(request.Password))
	if err != nil {
		err = errors.New("invalid username or password")
		return
	}

	role := "admin"
	if userM.IsSuperAdmin {
		role = "root"
	}
	accessToken, err := u.GenerateAccessToken(userM.ID, role)
	if err != nil {
		return
	}

	refreshTokenM, err := u.GenerateRefreshToken(userM.ID, nil)
	if err != nil {
		return
	}

	res = response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenM,
	}
	return
}
