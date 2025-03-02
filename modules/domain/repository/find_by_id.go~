package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(domainID uint64, ctx context.Context) (domainM entity.Domain, err error) {

	err = r.db.
		WithContext(ctx).
		Where("id = ?", domainID).
		Take(&domainM).Error
	if err != nil {
		err = errors.New("error find domain by id")
		return
	}

	return
}
