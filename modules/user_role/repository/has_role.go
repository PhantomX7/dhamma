package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) HasRole(ctx context.Context, userID, roleID uint64) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx)

	err := db.Model(&entity.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error

	if err != nil {
		return false, errors.WrapError(errors.ErrDatabase, "error checking role association")
	}

	return count > 0, nil
}
