package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/constants"
	"gorm.io/gorm"
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

	var accessToken, refreshToken string
	err = u.transactionManager.ExecuteInTransaction(func(tx *gorm.DB) error {
		err = u.userRepo.Create(ctx, &user, tx)
		if err != nil {
			return err
		}

		role := constants.EnumRoleAdmin
		if user.IsSuperAdmin {
			role = constants.EnumRoleRoot
		}

		accessToken, err = u.GenerateAccessToken(user.ID, role)
		if err != nil {
			return err
		}

		refreshToken, err = u.GenerateRefreshToken(user.ID, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	res = response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}
