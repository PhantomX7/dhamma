package repository

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

func (r *repository) AssignMultipleRole(userID uint64, roleAssignments []struct {
	DomainID uint64
	RoleID   uint64
}, tx *gorm.DB, ctx context.Context) (err error) {
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
			err = errors.New("error assign multiple role")
			return
		}
	}

	return
}
