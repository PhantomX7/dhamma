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
	// Add other repository methods if needed
}

type Service interface {
	Index(ctx context.Context, paginationConfig *pagination.Pagination) ([]entity.Event, utility.PaginationMeta, error)
	Show(ctx context.Context, eventID uint64) (entity.Event, error)
	Create(ctx context.Context, request request.EventCreateRequest) (entity.Event, error)
	Update(ctx context.Context, eventID uint64, request request.EventUpdateRequest) (entity.Event, error)
	Attend(ctx context.Context, eventID uint64, req request.EventAttendRequest) (entity.EventAttendance, error)
	AttendById(ctx context.Context, eventID uint64, req request.EventAttendByIDRequest) (entity.EventAttendance, error)
}

type Controller interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Attend(c *gin.Context)
	AttendById(c *gin.Context)
}
