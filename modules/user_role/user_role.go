package user_role

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	AssignRole(userID, domainID, roleID uint64, tx *gorm.DB, ctx context.Context) error
}
