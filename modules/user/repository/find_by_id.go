package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByID(userID uint64, ctx context.Context) (userM entity.User, err error) {

	err = r.db.Where("id = ?", userID).Take(&userM).Error
	if err != nil {
		err = errors.New("cannot find user with id")
		return
	}

	return
}
