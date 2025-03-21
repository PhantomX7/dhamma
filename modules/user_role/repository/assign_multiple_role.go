package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"gorm.io/gorm"
)

func (r *repository) AssignMultipleRole(ctx context.Context, userID uint64, roleAssignments []struct {
	DomainID uint64
	RoleID   uint64
}, tx *gorm.DB) (err error) {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	for _, ra := range roleAssignments {
		userRole := entity.UserRole{
			UserID:   userID,
			DomainID: ra.DomainID,
			RoleID:   ra.RoleID,
		}
		if err = tx.WithContext(ctx).Create(&userRole).Error; err != nil {
			err = utility.LogError("error assign multiple role", err)
			return
		}
	}

	return
}
