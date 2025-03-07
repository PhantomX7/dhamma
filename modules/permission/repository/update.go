package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) Update(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).Save(permission).Error
	if err != nil {
		return utility.LogError("error update permission", err)
	}
	return nil
}
