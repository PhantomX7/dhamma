package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) HasDomain(ctx context.Context, userID, domainID uint64) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx)

	err := db.Model(&entity.UserDomain{}).
		Where("user_id = ? AND domain_id = ?", userID, domainID).
		Count(&count).Error

	if err != nil {
		return false, utility.WrapError(utility.ErrDatabase, "error checking domain association")
	}

	return count > 0, nil
}
