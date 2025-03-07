package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByCodes(ctx context.Context, permissionCodes []string) (permissions []entity.Permission, err error) {

	err = r.db.WithContext(ctx).
		Where("code IN ?", permissionCodes).Find(&permissions).Error
	if err != nil {
		err = utility.LogError("error find permission by code", err)
		return
	}

	return
}
