package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

func (r *repository) AssignRole(userID, domainID, roleID uint64, tx *gorm.DB, ctx context.Context) (err error) {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	userRole := entity.UserRole{
		UserID:   userID,
		DomainID: domainID,
		RoleID:   roleID,
	}

	err = tx.WithContext(ctx).
		Create(&userRole).Error
	if err != nil {
		return errors.New("error assign role")
	}

	return
}
