package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) FindByUserIDAndDomainID(ctx context.Context, userID uint64, domainID uint64, preloadRelations bool) ([]entity.UserRole, error) {
	db := r.db.WithContext(ctx)

	if preloadRelations {
		db = db.Preload("Domain").Preload("Role")
	}

	var userRoles []entity.UserRole
	err := db.Where("user_id = ? AND domain_id = ?", userID, domainID).Find(&userRoles).Error

	if err != nil {
		return nil, utility.WrapError(utility.ErrDatabase, "error finding roles for user in domain")
	}

	return userRoles, nil
}
