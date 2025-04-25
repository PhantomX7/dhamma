package repository

import (
	"context"
	"time"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) FindByID(ctx context.Context, refreshTokenID string) (entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	db := r.db.WithContext(ctx)

	result := db.Where("id = ? AND is_valid = ? AND expires_at > ?",
		refreshTokenID, true, time.Now()).
		Take(&refreshToken)

	if result.Error != nil {
		return refreshToken, errors.WrapError(errors.ErrNotFound, "refresh token not found or expired")
	}

	return refreshToken, nil
}
