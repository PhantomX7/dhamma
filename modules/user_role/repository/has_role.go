package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) HasRole(ctx context.Context, userID, roleID uint64) (bool, error) {
	exist, err := r.base.Exists(ctx, map[string]any{
		"user_id": userID,
		"role_id": roleID,
	})
	if err != nil {
		return false, errors.WrapError(errors.ErrDatabase, "error checking role association")
	}

	return exist, nil
}
