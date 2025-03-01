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
	Create(user *entity.User, tx *gorm.DB, ctx context.Context) error
	Update(user *entity.User, tx *gorm.DB, ctx context.Context) error
	FindAll(pg *pagination.Pagination, ctx context.Context) ([]entity.User, error)
	FindByID(userID uint64, preloadRelations bool, ctx context.Context) (entity.User, error)
	FindByUsername(username string, ctx context.Context) (entity.User, error)
	FindByIDWithRelation(userID uint64, ctx context.Context) (entity.User, error)
	GetUserDomains(userID uint64, ctx context.Context) (entity.User, error)
	Count(pg *pagination.Pagination, ctx context.Context) (int64, error)
}

type Service interface {
	Index(pg *pagination.Pagination, ctx context.Context) ([]entity.User, utility.PaginationMeta, error)
	Show(userID uint64, ctx context.Context) (entity.User, error)
	Create(request request.UserCreateRequest, ctx context.Context) (entity.User, error)
	AssignDomain(userID uint64, request request.AssignDomainRequest, ctx context.Context) error
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Create(ctx *gin.Context)
	AssignDomain(ctx *gin.Context)
}
