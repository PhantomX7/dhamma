package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) HasDomain(ctx context.Context, userID, domainID uint64) (bool, error) {
	exist, err := r.base.Exists(ctx, map[string]any{
		"user_id":   userID,
		"domain_id": domainID,
	})

	if err != nil {
		return false, errors.WrapError(errors.ErrDatabase, "error checking domain association")
	}

	return exist, nil
}
