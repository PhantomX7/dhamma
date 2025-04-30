package user

import (
	"context"

	"github.com/PhantomX7/dhamma/modules/user/dto/request"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.User]
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.User, utility.PaginationMeta, error)
	Show(ctx context.Context, userID uint64) (entity.User, error)
	Create(ctx context.Context, request request.UserCreateRequest) (entity.User, error)
	AssignDomain(ctx context.Context, userID uint64, request request.AssignDomainRequest) error
	AssignRole(ctx context.Context, userID uint64, request request.AssignRoleRequest) error
	RemoveDomain(ctx context.Context, userID uint64, request request.RemoveDomainRequest) error
	RemoveRole(ctx context.Context, userID uint64, request request.RemoveRoleRequest) error
}

// Update the Controller interface to include RemoveDomain and RemoveRole
type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Create(ctx *gin.Context)
	AssignDomain(ctx *gin.Context)
	AssignRole(ctx *gin.Context)
	RemoveDomain(ctx *gin.Context)
	RemoveRole(ctx *gin.Context)
}
