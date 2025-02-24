package repository

import (
	"context"
	"errors"
	"time"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(refreshTokenID string, ctx context.Context) (refreshToken entity.RefreshToken, err error) {
	err = r.db.
		Where("user_id = ? AND is_valid = ? AND expired_at > ?",
			refreshTokenID, true, time.Now(),
		).
		Take(&refreshToken).Error
	if err != nil {
		err = errors.New("error find refresh token by id")
		return
	}

	return
}
