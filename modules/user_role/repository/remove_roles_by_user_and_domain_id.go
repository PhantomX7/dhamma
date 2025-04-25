package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"

	"gorm.io/gorm"
)

func (r *repository) RemoveRolesByUserAndDomainID(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error {
	db := r.prepareDB(ctx, tx)

	result := db.
		Where("user_id = ? AND domain_id = ?", userID, domainID).
		Delete(&entity.UserRole{})

	if result.Error != nil {
		return errors.WrapError(errors.ErrDatabase, "error removing roles for domain")
	}

	return nil
}
