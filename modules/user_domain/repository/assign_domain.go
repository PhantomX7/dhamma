package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) AssignDomain(ctx context.Context, userID, domainID uint64, tx *gorm.DB) error {
	userDomain := entity.UserDomain{
		UserID:   userID,
		DomainID: domainID,
	}

	return r.base.Create(ctx, &userDomain, tx)
}
