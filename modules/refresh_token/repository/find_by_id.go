package repository

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	customErrors "github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
)

func (r *repository) FindByID(ctx context.Context, refreshTokenID string) (entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	db := r.db.WithContext(ctx)

	result := db.Where("id = ? AND is_valid = ? AND expires_at > ?",
		refreshTokenID, true, time.Now()).
		Take(&refreshToken)

	if result.Error != nil {
		errMessage := "refresh token not found"
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Return not found error
			return refreshToken, customErrors.ErrNotFound
		}

		// Log the specific error
		logger.FromCtx(ctx).Error(errMessage, zap.String("id", refreshTokenID), zap.Error(result.Error))
		return refreshToken, customErrors.WrapError(customErrors.ErrNotFound, "refresh token not found or expired")
	}

	return refreshToken, nil
}
