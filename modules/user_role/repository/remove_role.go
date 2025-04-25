package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"

	"gorm.io/gorm"
)

func (r *repository) RemoveRole(ctx context.Context, userID, roleID uint64, tx *gorm.DB) error {
	// For delete operations with conditions, we need to use the DB directly
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&entity.UserRole{}).Error

	if err != nil {
		return errors.WrapError(errors.ErrDatabase, "error removing role")
	}

	return nil
}
