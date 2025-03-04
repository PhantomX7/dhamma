package role

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, role *entity.Role, tx *gorm.DB) error
	Update(ctx context.Context, role *entity.Role, tx *gorm.DB) error
	FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Role, error)
	FindByID(ctx context.Context, roleID uint64) (entity.Role, error)
	FindByNameAndDomainID(ctx context.Context, name string, domainID uint64) (entity.Role, error)
	Count(ctx context.Context, pg *pagination.Pagination) (int64, error)
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Role, utility.PaginationMeta, error)
	Show(ctx context.Context, roleID uint64) (entity.Role, error)
	Update(ctx context.Context, roleID uint64, request request.RoleUpdateRequest) (entity.Role, error)
	Create(ctx context.Context, request request.RoleCreateRequest) (entity.Role, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
