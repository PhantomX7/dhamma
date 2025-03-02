package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) InvalidateAllByUserID(userID uint64, ctx context.Context) (err error) {
	err = r.db.
		WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_valid", false).Error
	if err != nil {
		err = errors.New("error invalidate all refresh token by user id")
		return
	}

	return
}
