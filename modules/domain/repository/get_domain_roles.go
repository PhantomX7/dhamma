package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) GetDomainRoles(ctx context.Context, domainID uint64) (domainM entity.Domain, err error) {

	err = r.db.
		WithContext(ctx).
		Where("id = ?", domainID).
		Take(&domainM).Error
	if err != nil {
		err = utility.LogError("error get domain roles", err)
		return
	}

	return
}
