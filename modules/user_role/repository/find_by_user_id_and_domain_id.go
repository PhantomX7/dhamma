package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"gorm.io/gorm"
)

func (r *repository) FindByUserIDAndDomainID(ctx context.Context, userID uint64, domainID uint64, preloadRelations bool) (userRoles []entity.UserRole, err error) {
	var preloadScope = func(db *gorm.DB) *gorm.DB { return db }
	if preloadRelations {
		preloadScope = func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Domain").
				Preload("Role")
		}
	}

	err = r.db.WithContext(ctx).
		Scopes(preloadScope).
		Where("user_id = ? AND domain_id = ?", userID, domainID).
		Find(&userRoles).Error
	if err != nil {
		err = utility.LogError("error get user role detail by user id and domain id", err)
		return
	}

	return
}
