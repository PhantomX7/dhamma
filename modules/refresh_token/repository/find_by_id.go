package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"
	"time"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(refreshTokenID string, ctx context.Context) (refreshToken entity.RefreshToken, err error) {
	err = r.db.
		WithContext(ctx).
		Where("id = ? AND is_valid = ? AND expires_at > ?",
			refreshTokenID, true, time.Now(),
		).
		Take(&refreshToken).Error
	if err != nil {
		err = utility.LogError("error find refresh token by id", err)
		return
	}

	return
}
