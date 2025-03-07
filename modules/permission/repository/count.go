package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	var count int64

	filterScopes, _ := pagination.NewScopeBuilder(pg).Build()

	err := r.db.
		WithContext(ctx).
		Model(&entity.Permission{}).
		Scopes(filterScopes...).
		Count(&count).Error
	if err != nil {
		err = utility.LogError("error count permissions", err)
		return 0, err
	}

	return count, nil
}
