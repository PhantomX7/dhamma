package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByUsername(username string, ctx context.Context) (userM entity.User, err error) {

	err = r.db.Where("username = ?", username).First(&userM).Error
	if err != nil {
		err = errors.New("cannot find user with username")
		return
	}

	return userM, nil
}
