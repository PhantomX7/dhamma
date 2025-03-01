package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByCode(code string, ctx context.Context) (domainM entity.Domain, err error) {
	err = r.db.
		WithContext(ctx).
		Where("code = ?", code).
		First(&domainM).
		Error
	if err != nil {
		err = errors.New("error find domain by code")
		return
	}

	return
}
