package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) FindByUserID(ctx context.Context, userID uint64, preloads ...string) ([]entity.UserRole, error) {
	userRoles, err := r.base.FindByField(ctx, "user_id", userID, preloads...)
	if err != nil {
		return nil, errors.WrapError(errors.ErrDatabase, "error finding roles for user")
	}

	return userRoles, nil
}
