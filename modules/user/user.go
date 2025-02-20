package user

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/model"
	"github.com/PhantomX7/go-core/utility/request_util"
)

type Repository interface {
	Insert(user *model.User, tx *gorm.DB) error
	Update(user *model.User, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]model.User, error)
	FindByID(userID uint64) (model.User, error)
	FindByUsername(username string) (model.User, error)
	Count(config request_util.PaginationConfig) (int64, error)
}
