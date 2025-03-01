package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

func (r *repository) FindAll(pg *pagination.Pagination, ctx context.Context) ([]entity.Domain, error) {
	results := make([]entity.Domain, 0)

	filterScopes, metaScopes := pagination.NewScopeBuilder(pg).Build()

	err := r.db.
		WithContext(ctx).
		Scopes(filterScopes...).
		Scopes(metaScopes...).
		Find(&results).Error
	if err != nil {
		err = utility.LogError("error find domains", err)
		return results, err
	}

	return results, nil
}
