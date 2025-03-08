package permission

import (
	"context"
	"github.com/PhantomX7/dhamma/modules/permission/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/pagination"

	"gorm.io/gorm"
)

type permission struct {
	Key   string
	Index string
}

var Permissions = permission{
	Key:   "role",
	Index: "index",
}

type Repository interface {
	Create(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error
	Update(ctx context.Context, permission *entity.Permission, tx *gorm.DB) error
	FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.Permission, error)
	FindByID(ctx context.Context, permissionID uint64) (entity.Permission, error)
	FindByCode(ctx context.Context, permissionCode string) (entity.Permission, error)
	FindByCodes(ctx context.Context, permissionCodes []string) ([]entity.Permission, error)
	Count(ctx context.Context, pg *pagination.Pagination) (int64, error)
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Permission, utility.PaginationMeta, error)
	Show(ctx context.Context, permissionID uint64) (entity.Permission, error)
	Update(ctx context.Context, permissionID uint64, request request.PermissionUpdateRequest) (entity.Permission, error)
	Create(ctx context.Context, request request.PermissionCreateRequest) (entity.Permission, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
