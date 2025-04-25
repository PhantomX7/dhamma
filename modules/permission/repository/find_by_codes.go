package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func (r *repository) FindByCodes(ctx context.Context, permissionCodes []string) ([]entity.Permission, error) {
	var permissions []entity.Permission
	db := r.db.WithContext(ctx)

	err := db.Where("code IN ?", permissionCodes).Find(&permissions).Error
	if err != nil {
		return nil, errors.WrapError(err, "failed to find permissions by codes")
	}

	return permissions, nil
}
