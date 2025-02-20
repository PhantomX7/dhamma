package service

import (
	"net/http"

	"github.com/PhantomX7/go-core/utility/errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (u *service) UpdatePassword(
	request request.UpdatePasswordRequest,
	ctx *gin.Context,
) (err error) {
	userM, err := u.userRepo.FindByID(utility.GetIDFromContext(ctx))
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userM.Password), []byte(request.CurrentPassword))
	if err != nil {
		err = errors.CustomError{
			Message:  "password salah",
			HTTPCode: http.StatusUnprocessableEntity,
		}
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	userM.Password = string(password)

	err = u.userRepo.Update(&userM, nil)
	if err != nil {
		return
	}

	return
}
