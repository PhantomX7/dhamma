package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/utility"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(ctx context.Context, userID uint64, preloadRelations bool) (userM entity.User, err error) {

	var preloadScope = func(db *gorm.DB) *gorm.DB { return db }
	if preloadRelations {
		preloadScope = func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Domains").
				Preload("UserRoles.Role")
		}
	}
	err = r.db.WithContext(ctx).
		Scopes(preloadScope).
		Where("id = ?", userID).
		Take(&userM).Error
	if err != nil {
		err = utility.LogError("error find user by id", err)
		return
	}

	return
}
