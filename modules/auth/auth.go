package auth

import (
	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignIn(request request.SignInRequest) (response.AuthResponse, error)
	SignUp(request request.SignUpRequest) (response.AuthResponse, error)
	UpdatePassword(request request.UpdatePasswordRequest, ctx *gin.Context) error
	GetMe(ctx *gin.Context) (response.MeResponse, error)
}

type Controller interface {
	GetMe(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	SignUp(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
}
