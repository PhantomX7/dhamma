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
	// Add other repository methods if needed
}

type Service interface {
	Index(ctx context.Context, paginationConfig *pagination.Pagination) ([]entity.Follower, utility.PaginationMeta, error)
	Show(ctx context.Context, followerID uint64) (entity.Follower, error)
	Create(ctx context.Context, request request.FollowerCreateRequest) (entity.Follower, error)
	Update(ctx context.Context, followerID uint64, request request.FollowerUpdateRequest) (entity.Follower, error)
	AddCard(ctx context.Context, followerID uint64, req request.FollowerAddCardRequest) (entity.Card, error)
	DeleteCard(ctx context.Context, followerID uint64, cardID uint64) error
}

type Controller interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	AddCard(c *gin.Context)
	DeleteCard(c *gin.Context)
}
