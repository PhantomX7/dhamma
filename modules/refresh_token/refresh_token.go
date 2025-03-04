package refresh_token

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error
	Update(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error
	Delete(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error
	FindByID(ctx context.Context, refreshTokenID string) (entity.RefreshToken, error)
	GetValidCountByUserID(ctx context.Context, userID uint64) (int64, error)
	DeleteInvalidToken(ctx context.Context) error
	InvalidateAllByUserID(ctx context.Context, userID uint64) error
}
