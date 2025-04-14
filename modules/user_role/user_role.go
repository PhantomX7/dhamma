package user_role

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

// Add the new method to the Repository interface
type Repository interface {
	AssignRole(ctx context.Context, userID, domainID, roleID uint64, tx *gorm.DB) error
	AssignMultipleRole(ctx context.Context, userID uint64, roleAssignments []struct {
		DomainID uint64
		RoleID   uint64
	}, tx *gorm.DB) error
	FindByUserID(ctx context.Context, userID uint64, preloadRelations bool) ([]entity.UserRole, error)
	FindByUserIDAndDomainID(ctx context.Context, userID uint64, domainID uint64, preloadRelations bool) ([]entity.UserRole, error)
	HasRole(ctx context.Context, userID, roleID uint64) (bool, error)
	RemoveRole(ctx context.Context, userID, roleID uint64, tx *gorm.DB) error
	RemoveRolesByUserAndDomainID(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error
}
