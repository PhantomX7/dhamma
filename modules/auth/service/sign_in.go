package service

import (
	"context"
	"strings"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/utility/errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

// SignIn handles root/admin authentication only
func (s *service) SignIn(ctx context.Context, request request.SignInRequest) (res response.AuthResponse, err error) {
	user := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	user, err = s.userRepo.FindOneByField(ctx, "username", request.Username)
	if err != nil {
		err = errors.NewServiceError("invalid username or password", nil)
		return
	}

	// Only allow super admin users to sign in through root route
	if !user.IsSuperAdmin {
		err = errors.NewServiceError("access denied: only root users can sign in through this route", nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		err = errors.NewServiceError("invalid username or password", nil)
		return
	}

	role := constants.EnumRoleRoot // Root users always get root role

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
