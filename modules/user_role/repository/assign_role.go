package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) AssignRole(ctx context.Context, userID, domainID, roleID uint64, tx *gorm.DB) error {
	userRole := entity.UserRole{
		UserID:   userID,
		DomainID: domainID,
		RoleID:   roleID,
	}

	return r.base.Create(ctx, &userRole, tx)
}
