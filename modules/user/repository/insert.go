package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/errors"
)

func (r *repository) Insert(user *entity.User, tx *gorm.DB) error {
	var db = r.db
	if tx != nil {
		db = tx
	}
	err := db.Create(user).Error
	if err != nil {
		log.Println("error-insert-user:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}
