package event_attendance

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event_attendance/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.EventAttendance]
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.EventAttendance, utility.PaginationMeta, error)
	Show(ctx context.Context, eventAttendanceID uint64) (entity.EventAttendance, error)
	Update(ctx context.Context, eventAttendanceID uint64, request request.EventAttendanceUpdateRequest) (entity.EventAttendance, error)
	Create(ctx context.Context, request request.EventAttendanceCreateRequest) (entity.EventAttendance, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
