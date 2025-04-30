package auth

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/modules/auth/dto/request"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignIn(ctx context.Context, request request.SignInRequest) (response.AuthResponse, error)
	SignUp(ctx context.Context, request request.SignUpRequest) (response.AuthResponse, error)
	Refresh(ctx context.Context, request request.RefreshRequest) (response.AuthResponse, error)
	UpdatePassword(ctx context.Context, request request.UpdatePasswordRequest) error
	GetMe(ctx context.Context) (response.MeResponse, error)
	GenerateAccessToken(userID uint64, role string) (string, error)
	GenerateRefreshToken(userID uint64, tx *gorm.DB) (string, error)
}

type Controller interface {
	GetMe(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	SignUp(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
}
