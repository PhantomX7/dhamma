package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

func (r *repository) Count(pg *pagination.Pagination, ctx context.Context) (int64, error) {
	var count int64

	filterScopes, _ := pagination.NewScopeBuilder(pg).Build()

	err := r.db.
		WithContext(ctx).
		Model(&entity.Domain{}).
		Scopes(filterScopes...).
		Count(&count).Error
	if err != nil {
		err = errors.New("error count domains")
		return 0, err
	}

	return count, nil
}
