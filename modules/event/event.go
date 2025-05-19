package event

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.Event]
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.Event, utility.PaginationMeta, error)
	Show(ctx context.Context, eventID uint64) (entity.Event, error)
	Update(ctx context.Context, eventID uint64, request request.EventUpdateRequest) (entity.Event, error)
	Create(ctx context.Context, request request.EventCreateRequest) (entity.Event, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
