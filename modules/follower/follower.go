package follower

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/follower/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.Follower]
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Follower, utility.PaginationMeta, error)
	Show(ctx context.Context, followerID uint64) (entity.Follower, error)
	Update(ctx context.Context, followerID uint64, request request.FollowerUpdateRequest) (entity.Follower, error)
	Create(ctx context.Context, request request.FollowerCreateRequest) (entity.Follower, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
