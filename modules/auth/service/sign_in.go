package service

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/go-core/utility/errors"
)

func (u *service) SignIn(request request.SignInRequest, ctx context.Context) (res response.AuthResponse, err error) {
	userM := entity.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	userM, err = u.userRepo.FindByUsername(request.Username, ctx)
	if err != nil {
		err = errors.ErrFailedAuthentication
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userM.Password), []byte(request.Password))
	if err != nil {
		err = errors.ErrFailedAuthentication
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		ID:       userM.ID,
		Username: userM.Username,
		IssuedAt: time.Now().Unix(),
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * 80000),
			},
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		err = errors.ErrFailedAuthentication
		return
	}

	res = response.AuthResponse{
		Token: tokenString,
	}
	return
}
