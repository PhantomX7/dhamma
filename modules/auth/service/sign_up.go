package service

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/model"
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/go-core/utility/errors"
)

func (u *service) SignUp(request request.SignUpRequest) (res response.AuthResponse, err error) {
	userM := model.User{}

	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	_ = copier.Copy(&userM, &request)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.ErrFailedAuthentication
		return
	}
	userM.Password = string(password)

	err = u.userRepo.Insert(&userM, nil)
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
