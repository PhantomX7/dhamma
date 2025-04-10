package role

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"
	"github.com/gin-gonic/gin"
)

type permission struct {
	Key            string
	Index          string
	Show           string
	Create         string
	Update         string
	AddPermissions string
}

var Permissions = permission{
	Key:            "role",
	Index:          "index",
	Show:           "show",
	Create:         "create",
	Update:         "update",
	AddPermissions: "add-permissions",
}

type Repository interface {
	repository.BaseRepositoryInterface[entity.Role]
	FindByNameAndDomainID(ctx context.Context, name string, domainID uint64) (entity.Role, error)
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Role, utility.PaginationMeta, error)
	Show(ctx context.Context, roleID uint64) (entity.Role, error)
	Update(ctx context.Context, roleID uint64, request request.RoleUpdateRequest) (entity.Role, error)
	Create(ctx context.Context, request request.RoleCreateRequest) (entity.Role, error)
	AddPermissions(ctx context.Context, roleID uint64, request request.RoleAddPermissionsRequest) (entity.Role, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
	AddPermissions(ctx *gin.Context)
}
