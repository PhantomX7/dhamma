package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) InvalidateAllByUserID(ctx context.Context, userID uint64) (err error) {
	err = r.db.
		WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_valid", false).Error
	if err != nil {
		err = utility.LogError("error invalidate all refresh token by user id", err)
		return
	}

	return
}
