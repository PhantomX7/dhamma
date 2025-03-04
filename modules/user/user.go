package user

import (
	"context"
	"github.com/PhantomX7/dhamma/modules/user/dto/request"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User, tx *gorm.DB) error
	Update(ctx context.Context, user *entity.User, tx *gorm.DB) error
	FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.User, error)
	FindByID(ctx context.Context, userID uint64, preloadRelations bool) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	FindByIDWithRelation(ctx context.Context, userID uint64) (entity.User, error)
	GetUserDomains(ctx context.Context, userID uint64) (entity.User, error)
	Count(ctx context.Context, pg *pagination.Pagination) (int64, error)
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.User, utility.PaginationMeta, error)
	Show(ctx context.Context, userID uint64) (entity.User, error)
	Create(ctx context.Context, request request.UserCreateRequest) (entity.User, error)
	AssignDomain(ctx context.Context, userID uint64, request request.AssignDomainRequest) error
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Create(ctx *gin.Context)
	AssignDomain(ctx *gin.Context)
}
