package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)


func (r *repository) Update(refreshToken *entity.RefreshToken, tx *gorm.DB, ctx context.Context) error {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).Save(refreshToken).Error
	if err != nil {
		return utility.LogError("error create refresh token", err)
	}
	return nil
}
