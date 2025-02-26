package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(roleID uint64, ctx context.Context) (roleM entity.Role, err error) {

	err = r.db.WithContext(ctx).
		Preload("Domain").
		Where("id = ?", roleID).Take(&roleM).Error
	if err != nil {
		err = errors.New("error find role by id")
		return
	}

	return
}
