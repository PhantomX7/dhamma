package repository

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) HasDomain(userID, domainID uint64, ctx context.Context) (bool bool, err error) {
	var count int64
	err = r.db.WithContext(ctx).
		Model(&entity.UserDomain{}).
		Where("user_id = ? AND domain_id = ?", userID, domainID).
		Count(&count).Error
	if err != nil {
		return false, errors.New("error check has domain")
	}

	return count > 0, nil
}
