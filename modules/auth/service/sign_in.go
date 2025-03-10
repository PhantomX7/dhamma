package service

import (
	"context"
	"errors"
	"strings"

	"github.com/PhantomX7/dhamma/constants"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (s *service) SignIn(ctx context.Context, request request.SignInRequest) (res response.AuthResponse, err error) {
	user := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	user, err = s.userRepo.FindByUsername(ctx, request.Username)
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

	accessToken, err := s.GenerateAccessToken(user.ID, role)
	if err != nil {
		return
	}

	refreshTokenM, err := s.GenerateRefreshToken(user.ID, nil)
	if err != nil {
		return
	}

	res = response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenM,
	}
	return
}
