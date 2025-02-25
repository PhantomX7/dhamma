package repository

import (
	"context"
	"errors"
	"time"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) DeleteInvalidToken(ctx context.Context) (err error) {
	err = r.db.
		WithContext(ctx).
		Where("expires_at < ? OR is_valid = ?", time.Now(), false).
		Delete(&entity.RefreshToken{}).Error
	if err != nil {
		err = errors.New("error count valid refresh token")
		return
	}

	return
}
