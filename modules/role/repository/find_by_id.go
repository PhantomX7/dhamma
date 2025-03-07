package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(ctx context.Context, roleID uint64) (roleM entity.Role, err error) {

	err = r.db.WithContext(ctx).
		Preload("Domain").
		Where("id = ?", roleID).Take(&roleM).Error
	if err != nil {
		err = utility.LogError("error find role by id", err)
		return
	}

	return
}
