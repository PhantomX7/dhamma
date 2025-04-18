package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) RemoveDomain(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error {
	db := r.prepareDB(ctx, tx)

	result := db.Where("user_id = ? AND domain_id = ?", userID, domainID).Delete(&entity.UserDomain{})
	if result.Error != nil {
		return utility.WrapError(utility.ErrDatabase, "failed to remove domain from user")
	}

	if result.RowsAffected == 0 {
		return utility.WrapError(utility.ErrNotFound, "domain assignment not found")
	}

	return nil
}

func (r *repository) prepareDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}
