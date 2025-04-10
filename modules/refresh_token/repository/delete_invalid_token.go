package repository

import (
	"context"
	"time"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) DeleteInvalidToken(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at < ? OR is_valid = ?", time.Now(), false).
		Delete(&entity.RefreshToken{}).Error

	if err != nil {
		return utility.WrapError(utility.ErrDatabase, "failed to delete invalid tokens")
	}

	return nil
}
