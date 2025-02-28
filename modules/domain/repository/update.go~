package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) Update(role *entity.Role, tx *gorm.DB, ctx context.Context) error {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).Save(role).Error
	if err != nil {
		return errors.New("error update role")
	}
	return nil
}
