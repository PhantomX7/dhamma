package event_attendance

import (
	"context"
	"time" // Add time import

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

// Repository defines the interface for event_attendance data operations.
type Repository interface {
	repository.BaseRepositoryInterface[entity.EventAttendance]
	// HasAttendedOnDate checks if a follower attended a specific event on a given date.
	HasAttendedOnDate(ctx context.Context, followerID uint64, eventID uint64, date time.Time) (bool, error)
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.EventAttendance, utility.PaginationMeta, error)
	Show(ctx context.Context, eventAttendanceID uint64) (entity.EventAttendance, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
}
