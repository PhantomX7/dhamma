package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(ctx context.Context, domainID uint64) (domainM entity.Domain, err error) {

	err = r.db.
		WithContext(ctx).
		Where("id = ?", domainID).
		Take(&domainM).Error
	if err != nil {
		err = utility.LogError("error find domain by id", err)
		return
	}

	return
}
