package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) AssignRole(ctx context.Context, userID, domainID, roleID uint64, tx *gorm.DB) (err error) {
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
		return utility.LogError("error assign role", err)
	}

	return
}
