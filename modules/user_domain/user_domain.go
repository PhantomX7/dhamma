package user_domain

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"gorm.io/gorm"
)

type Repository interface {
	// repository.BaseRepositoryInterface[entity.UserDomain]
	HasDomain(ctx context.Context, userID, domainID uint64) (bool, error)
	AssignDomain(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error
	RemoveDomain(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error
	FindByUserID(ctx context.Context, userID uint64, preloads ...string) ([]entity.UserDomain, error)
}
