package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) AssignDomain(userID, domainID uint64, tx *gorm.DB, ctx context.Context) (err error) {
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
		return errors.New("error assign role")
	}

	return
}
