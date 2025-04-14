package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"

	"gorm.io/gorm"
)

func (r *repository) RemoveRolesByUserAndDomainID(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).
		Where("user_id = ? AND domain_id = ?", userID, domainID).
		Delete(&entity.UserRole{})

	if result.Error != nil {
		return utility.WrapError(utility.ErrDatabase, "error removing roles for domain")
	}

	return nil
}
