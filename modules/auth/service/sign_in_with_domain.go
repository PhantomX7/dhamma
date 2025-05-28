package service

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// SignInWithDomain handles domain-specific authentication
func (s *service) SignInWithDomain(ctx context.Context, request request.SignInRequest, domainCode string) (res response.AuthResponse, err error) {
	user := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	user, err = s.userRepo.FindOneByField(ctx, "username", request.Username)
	if err != nil {
		err = errors.NewServiceError("invalid username or password", nil)
		return
	}

	// Root users cannot sign in through domain routes
	if user.IsSuperAdmin {
		err = errors.NewServiceError("access denied: root users must sign in through admin route", nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		err = errors.NewServiceError("invalid username or password", nil)
		return
	}

	// Validate that the domain exists
	domain, err := s.domainRepo.FindOneByField(ctx, "code", domainCode)
	if err != nil {
		err = errors.NewServiceError("invalid domain", nil)
		return
	}

	// Check if user has access to this domain
	hasAccess, err := s.userDomainRepo.HasDomain(ctx, user.ID, domain.ID)
	if err != nil {
		err = errors.NewServiceError("error checking domain access", err)
		return
	}

	if !hasAccess {
		err = errors.NewServiceError("access denied: user does not have access to this domain", nil)
		return
	}

	role := constants.EnumRoleAdmin // Domain users get admin role

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
