package role

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(role *entity.Role, tx *gorm.DB, ctx context.Context) error
	Update(role *entity.Role, tx *gorm.DB, ctx context.Context) error
	FindAll(pg *pagination.Pagination, ctx context.Context) ([]entity.Role, error)
	FindByID(roleID uint64, ctx context.Context) (entity.Role, error)
	FindByName(name string, ctx context.Context) (entity.Role, error)
	Count(pg *pagination.Pagination, ctx context.Context) (int64, error)
}
