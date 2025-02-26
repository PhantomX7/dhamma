package repository

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) GetUserDomains(userID uint64, ctx context.Context) (userM entity.User, err error) {

	err = r.db.WithContext(ctx).
		Preload("Domains").
		Where("id = ?", userID).Take(&userM).Error
	if err != nil {
		err = errors.New("error get user domains")
		return
	}

	return
}
