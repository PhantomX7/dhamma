package repository

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) GetUserRoleDetailByUserID(userID uint64, ctx context.Context) (userRoles []entity.UserRole, err error) {
	err = r.db.WithContext(ctx).
		Preload("Domain").
		Preload("Role").
		Where("user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		err = errors.New("error get user role detail by user id")
		return
	}

	return
}
