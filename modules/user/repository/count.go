package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/request_util"
)

func (r *repository) Count(config request_util.PaginationConfig, ctx context.Context) (int64, error) {
	var count int64

	err := r.db.
		Model(&entity.User{}).
		Scopes(config.Scopes()...).
		Count(&count).Error
	if err != nil {
		err = errors.New("cannot count users")
		return 0, err
	}

	return count, nil
}
