package point_mutation

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.PointMutation]
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.PointMutation, utility.PaginationMeta, error)
	Show(ctx context.Context, pointMutationID uint64) (entity.PointMutation, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
}
