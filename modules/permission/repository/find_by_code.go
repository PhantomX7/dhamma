package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByCode(ctx context.Context, permissionCode string) (permissionM entity.Permission, err error) {

	err = r.db.WithContext(ctx).
		Where("code = ?", permissionCode).Take(&permissionM).Error
	if err != nil {
		err = utility.LogError("error find permission by code", err)
		return
	}

	return
}
