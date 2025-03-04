package user_domain

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

type Repository interface {
	HasDomain(ctx context.Context, userID, domainID uint64) (bool, error)
	AssignDomain(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error
	FindByUserID(ctx context.Context, userID uint64, preloadRelations bool) ([]entity.UserDomain, error)
}
