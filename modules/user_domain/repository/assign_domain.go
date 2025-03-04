package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) AssignDomain(ctx context.Context, userID, domainID uint64, tx *gorm.DB) (err error) {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	userDomain := entity.UserDomain{
		UserID:   userID,
		DomainID: domainID,
	}

	err = tx.WithContext(ctx).
		Create(&userDomain).Error
	if err != nil {
		return utility.LogError("error assign domain", err)
	}

	return
}
