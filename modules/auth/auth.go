package auth

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignIn(request request.SignInRequest, ctx context.Context) (response.AuthResponse, error)
	SignUp(request request.SignUpRequest, ctx context.Context) (response.AuthResponse, error)
	UpdatePassword(request request.UpdatePasswordRequest, ctx context.Context) error
	GetMe(userID uint64, ctx context.Context) (response.MeResponse, error)
}

type Controller interface {
	GetMe(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	SignUp(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
}
