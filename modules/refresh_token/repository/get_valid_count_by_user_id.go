package repository

import (
	"context"
	"time"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) GetValidCountByUserID(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_valid = ? AND expires_at > ?", userID, true, time.Now()).
		Count(&count).Error

	if err != nil {
		return 0, errors.WrapError(errors.ErrDatabase, "failed to count valid refresh tokens")
	}

	return count, nil
}
