package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByUsername(username string, ctx context.Context) (userM entity.User, err error) {

	err = r.db.WithContext(ctx).Where("username = ?", username).First(&userM).Error
	if err != nil {
		err = errors.New("error find user by username")
		return
	}

	return
}
