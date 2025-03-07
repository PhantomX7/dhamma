package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(ctx context.Context, permissionID uint64) (permissionM entity.Permission, err error) {

	err = r.db.WithContext(ctx).
		Where("id = ?", permissionID).Take(&permissionM).Error
	if err != nil {
		err = utility.LogError("error find permission by id", err)
		return
	}

	return
}
