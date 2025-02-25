package refresh_token

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(refreshToken *entity.RefreshToken, tx *gorm.DB, ctx context.Context) error
	Update(refreshToken *entity.RefreshToken, tx *gorm.DB, ctx context.Context) error
	Delete(refreshToken *entity.RefreshToken, tx *gorm.DB, ctx context.Context) error
	FindByID(refreshTokenID string, ctx context.Context) (entity.RefreshToken, error)
	GetValidCountByUserID(userID uint64, ctx context.Context) (int64, error)
	DeleteInvalidToken(ctx context.Context) error
	InvalidateAllByUserID(userID uint64, ctx context.Context) error
}
