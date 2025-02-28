package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/constants"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (u *service) SignIn(request request.SignInRequest, ctx context.Context) (res response.AuthResponse, err error) {
	user := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	user, err = u.userRepo.FindByUsername(request.Username, ctx)
	if err != nil {
		err = errors.New("invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		err = errors.New("invalid username or password")
		return
	}

	role := constants.EnumRoleAdmin
	if user.IsSuperAdmin {
		role = constants.EnumRoleRoot
	}

	accessToken, err := u.GenerateAccessToken(user.ID, role)
	if err != nil {
		return
	}

	refreshTokenM, err := u.GenerateRefreshToken(user.ID, nil)
	if err != nil {
		return
	}

	res = response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenM,
	}
	return
}
