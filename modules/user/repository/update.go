package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/model"
	"github.com/PhantomX7/go-core/utility/errors"
)

func (r *repository) Update(user *model.User, tx *gorm.DB) error {
	var db = r.db
	if tx != nil {
		db = tx
	}
	err := db.Save(user).Error
	if err != nil {
		log.Println("error-update-user:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}
