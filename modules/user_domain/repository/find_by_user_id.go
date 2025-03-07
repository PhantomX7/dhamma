package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) FindByUserID(ctx context.Context, userID uint64, preloadRelations bool) (userDomains []entity.UserDomain, err error) {
	var preloadScope = func(db *gorm.DB) *gorm.DB { return db }
	if preloadRelations {
		preloadScope = func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Domain")
		}
	}

	err = r.db.WithContext(ctx).
		Scopes(preloadScope).
		Where("user_id = ?", userID).
		Find(&userDomains).Error
	if err != nil {
		err = utility.LogError("error get user domain detail by user id", err)
		return
	}

	return
}
