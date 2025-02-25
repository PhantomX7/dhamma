package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) Create(role *entity.Role, tx *gorm.DB, ctx context.Context) error {
	// if tx is nil, use default db
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).Create(role).Error
	if err != nil {
		return errors.New("error create role")
	}
	return nil
}
