package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) GetUserDomains(userID uint64, ctx context.Context) (userM entity.User, err error) {

	err = r.db.WithContext(ctx).
		Preload("Domains").
		Where("id = ?", userID).Take(&userM).Error
	if err != nil {
		err = utility.LogError("error get user domains", err)
		return
	}

	return
}
