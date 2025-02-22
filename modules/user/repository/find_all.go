package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/request_util"
)

func (r *repository) FindAll(config request_util.PaginationConfig, ctx context.Context) ([]entity.User, error) {
	results := make([]entity.User, 0)

	err := r.db.
		Scopes(config.MetaScopes()...).
		Scopes(config.Scopes()...).
		Find(&results).Error
	if err != nil {
		err = errors.New("cannot find users")
		return results, err
	}

	return results, nil
}
