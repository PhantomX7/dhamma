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
	Create(role *entity.Role, tx *gorm.DB, ctx context.Context) error
	Update(role *entity.Role, tx *gorm.DB, ctx context.Context) error
	FindAll(pg *pagination.Pagination, ctx context.Context) ([]entity.Role, error)
	FindByID(roleID uint64, ctx context.Context) (entity.Role, error)
	FindByNameAndDomainID(name string, domainID uint64, ctx context.Context) (entity.Role, error)
	Count(pg *pagination.Pagination, ctx context.Context) (int64, error)
}

type Service interface {
	Index(pg *pagination.Pagination, ctx context.Context) ([]entity.Role, utility.PaginationMeta, error)
	Show(roleID uint64, ctx context.Context) (entity.Role, error)
	Update(roleID uint64, request request.RoleUpdateRequest, ctx context.Context) (entity.Role, error)
	Create(request request.RoleCreateRequest, ctx context.Context) (entity.Role, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
