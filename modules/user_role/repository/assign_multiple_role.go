package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) AssignMultipleRole(ctx context.Context, userID uint64, roleAssignments []struct {
	DomainID uint64
	RoleID   uint64
}, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	for _, ra := range roleAssignments {
		userRole := entity.UserRole{
			UserID:   userID,
			DomainID: ra.DomainID,
			RoleID:   ra.RoleID,
		}

		if err := r.base.Create(ctx, &userRole, db); err != nil {
			return errors.WrapError(errors.ErrDatabase, "error assigning multiple roles")
		}
	}

	return nil
}
