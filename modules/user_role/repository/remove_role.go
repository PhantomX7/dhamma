package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"

	"gorm.io/gorm"
)

func (r *repository) RemoveRole(userID, domainID, roleID uint64, tx *gorm.DB, ctx context.Context) (err error) {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	err = tx.WithContext(ctx).
		Where(
			"user_id = ? AND domain_id = ? AND role_id = ?",
			userID, domainID, roleID,
		).Delete(&entity.UserRole{}).Error
	if err != nil {
		return utility.LogError("error remove role", err)
	}

	return
}
