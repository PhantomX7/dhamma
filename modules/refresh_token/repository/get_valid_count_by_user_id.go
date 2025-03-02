package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"
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
		err = utility.LogError("error count valid refresh token", err)
		return
	}

	return
}
