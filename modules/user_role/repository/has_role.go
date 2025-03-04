package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) HasRole(ctx context.Context, userID, domainID, roleID uint64) (bool bool, err error) {
	var count int64
	err = r.db.WithContext(ctx).
		Model(&entity.UserRole{}).
		Where("user_id = ? AND domain_id = ? AND role_id = ?", userID, domainID, roleID).
		Count(&count).Error
	if err != nil {
		return false, utility.LogError("error check has role", err)
	}

	return count > 0, nil
}
