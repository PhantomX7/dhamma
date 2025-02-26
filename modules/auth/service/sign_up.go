package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
)

func (u *service) SignUp(request request.SignUpRequest, ctx context.Context) (res response.AuthResponse, err error) {
	userM := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	_ = copier.Copy(&userM, &request)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("failed to hash password")
		return
	}
	userM.Password = string(password)

	tx := u.transactionManager.NewTransaction()

	err = u.userRepo.Create(&userM, tx, ctx)
	if err != nil {
		tx.Rollback()
		return
	}

	accessToken, err := u.GenerateAccessToken(userM.ID, "admin")
	if err != nil {
		tx.Rollback()
		return
	}

	refreshToken, err := u.GenerateRefreshToken(userM.ID, nil)
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
