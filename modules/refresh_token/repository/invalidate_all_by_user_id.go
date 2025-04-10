package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) InvalidateAllByUserID(ctx context.Context, userID uint64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_valid", false).Error

	if err != nil {
		return utility.WrapError(utility.ErrDatabase, "failed to invalidate tokens for user")
	}

	return nil
}
