package repository

import (
	"log"

	"github.com/PhantomX7/dhamma/model"
	"github.com/PhantomX7/go-core/utility/errors"
	"github.com/PhantomX7/go-core/utility/request_util"
)

func (r *repository) FindAll(config request_util.PaginationConfig) ([]model.User, error) {
	results := make([]model.User, 0)

	err := r.db.
		Scopes(config.MetaScopes()...).
		Scopes(config.Scopes()...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-user:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}
