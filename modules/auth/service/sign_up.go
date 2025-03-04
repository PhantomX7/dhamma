package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/constants"
	"strings"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (u *service) SignUp(ctx context.Context, request request.SignUpRequest) (res response.AuthResponse, err error) {
	user := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	_ = copier.Copy(&user, &request)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("failed to hash password")
		return
	}
	user.Password = string(password)

	tx := u.transactionManager.NewTransaction()

	err = u.userRepo.Create(ctx, &user, tx)
	if err != nil {
		tx.Rollback()
		return
	}

	role := constants.EnumRoleAdmin
	if user.IsSuperAdmin {
		role = constants.EnumRoleRoot
	}

	accessToken, err := u.GenerateAccessToken(user.ID, role)
	if err != nil {
		tx.Rollback()
		return
	}

	refreshToken, err := u.GenerateRefreshToken(user.ID, nil)
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
