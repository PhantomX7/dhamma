package repository

import (
	"context"
	"errors"
	"time"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) GetValidCountByUserID(userID uint64, ctx context.Context) (count int64, err error) {
	err = r.db.
		WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_valid = ? AND expires_at > ?", userID, true, time.Now()).
		Count(&count).Error
	if err != nil {
		err = errors.New("error count valid refresh token")
		return
	}

	return
}
