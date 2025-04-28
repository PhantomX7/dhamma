package role

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request" // Add request DTO import
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"
	"github.com/gin-gonic/gin"
)

type permission struct {
	Key               string
	Index             string
	Show              string
	Create            string
	Update            string
	AddPermissions    string
	DeletePermissions string
}

var Permissions = permission{
	Key:               "role",
	Index:             "index",
	Show:              "show",
	Create:            "create",
	Update:            "update",
	AddPermissions:    "add-permissions",
	DeletePermissions: "delete-permissions",
}

type Repository interface {
	repository.BaseRepositoryInterface[entity.Role]
}

// Service defines the interface for role-related business logic.
type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Role, utility.PaginationMeta, error)
	Show(ctx context.Context, roleID uint64) (entity.Role, error)
	Update(ctx context.Context, roleID uint64, request request.RoleUpdateRequest) (entity.Role, error)
	Create(ctx context.Context, request request.RoleCreateRequest) (entity.Role, error)
	AddPermissions(ctx context.Context, roleID uint64, request request.RoleAddPermissionsRequest) (entity.Role, error)
	DeletePermissions(ctx context.Context, roleID uint64, request request.RoleDeletePermissionsRequest) error // Add DeletePermissions method
}

// Controller defines the interface for handling role-related HTTP requests.
type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
	AddPermissions(ctx *gin.Context)
	DeletePermissions(ctx *gin.Context) // Add DeletePermissions method
}
