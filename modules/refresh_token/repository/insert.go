package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) Create(ctx context.Context, refreshToken *entity.RefreshToken, tx *gorm.DB) error {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).Create(refreshToken).Error
	if err != nil {
		return utility.LogError("error create refresh token", err)
	}
	return nil
}
