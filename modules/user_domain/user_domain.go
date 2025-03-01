package user_domain

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	HasDomain(userID, domainID uint64, ctx context.Context) (bool, error)
	AssignDomain(userID, domainID uint64, tx *gorm.DB, ctx context.Context) error
}
