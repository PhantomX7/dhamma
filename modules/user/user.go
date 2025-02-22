package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/go-core/utility/request_util"
)

type Repository interface {
	Insert(user *entity.User, tx *gorm.DB, ctx context.Context) error
	Update(user *entity.User, tx *gorm.DB, ctx context.Context) error
	FindAll(config request_util.PaginationConfig, ctx context.Context) ([]entity.User, error)
	FindByID(userID uint64, ctx context.Context) (entity.User, error)
	FindByUsername(username string, ctx context.Context) (entity.User, error)
	Count(config request_util.PaginationConfig, ctx context.Context) (int64, error)
}
