package repository

import (
	"log"

	"github.com/PhantomX7/dhamma/model"
	"github.com/PhantomX7/go-core/utility/errors"
	"github.com/PhantomX7/go-core/utility/request_util"
)

func (r *repository) Count(config request_util.PaginationConfig) (int64, error) {
	var count int64

	err := r.db.
		Model(&model.User{}).
		Scopes(config.Scopes()...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-user:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
