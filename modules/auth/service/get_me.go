package service

import (
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"
)

func (u *service) GetMe(ctx *gin.Context) (res response.MeResponse, err error) {
	userM, err := u.userRepo.FindByID(utility.GetIDFromContext(ctx))
	if err != nil {
		return
	}

	res.User = userM

	return
}
