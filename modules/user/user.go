package user

import (
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/request_util"
)

type Repository interface {
	Insert(user *entity.User, tx *gorm.DB) error
	Update(user *entity.User, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]entity.User, error)
	FindByID(userID uint64) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
	Count(config request_util.PaginationConfig) (int64, error)
}
