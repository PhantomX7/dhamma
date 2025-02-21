package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/errors"
)

func (r *repository) FindByID(userID uint64) (userM entity.User, err error) {

	err = r.db.Where("id = ?", userID).Take(&userM).Error

	if err == gorm.ErrRecordNotFound {
		err = errors.ErrNotFound
		return
	}

	if err != nil {
		log.Println("error-find-user-by-id:", err)
		err = errors.ErrUnprocessableEntity
		return
	}

	return
}
