package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByName(name string, ctx context.Context) (roleM entity.Role, err error) {

	err = r.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&roleM).
		Error
	if err != nil {
		err = errors.New("error find role by name")
		return
	}

	return
}
