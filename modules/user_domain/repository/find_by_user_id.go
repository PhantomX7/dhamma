package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) FindByUserID(ctx context.Context, userID uint64, preloadRelations bool) ([]entity.UserDomain, error) {
	var preloads []string
	if preloadRelations {
		preloads = append(preloads, "Domain")
	}

	domains, err := r.base.FindByField(ctx, "user_id", userID, preloads...)
	if err != nil {
		return nil, utility.WrapError(utility.ErrDatabase, "error finding domains for user")
	}

	return domains, nil
}
