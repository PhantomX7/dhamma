package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/model"
	"github.com/PhantomX7/go-core/utility/errors"
)

func (r *repository) FindByUsername(username string) (userM model.User, err error) {

	err = r.db.Where("username = ?", username).First(&userM).Error

	if err == gorm.ErrRecordNotFound {
		err = errors.ErrNotFound
		return
	}

	if err != nil {
		log.Println("error-find-user-by-username:", err)
		err = errors.ErrUnprocessableEntity
		return
	}

	return userM, nil
}
