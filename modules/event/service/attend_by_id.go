package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/event/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"gorm.io/gorm"
)

// Attend handles the logic for a follower attending an event and receiving points.
func (s *service) AttendById(ctx context.Context, eventID uint64, req request.EventAttendByIDRequest) (eventAttendance entity.EventAttendance, err error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return
	}

	_, err = utility.CheckDomainContext(ctx, event.DomainID, "event", "attend")
	if err != nil {
		return
	}

	follower, err := s.followerRepo.FindByID(ctx, req.FollowerID)
	if err != nil {
		return
	}

	// 4. Check if follower has already attended this event
	atendded, err := s.eventAttendanceRepo.HasAttendedOnDate(
		ctx,
		follower.ID,
		event.ID,
		time.Now(),
	)
	if err != nil {
		return
	}

	if atendded {
		return eventAttendance, &errors.AppError{
			Message: fmt.Sprintf("follower has already attended the event"),
			Status:  http.StatusForbidden,
		}
	}

	if follower.DomainID != event.DomainID {
		return eventAttendance, &errors.AppError{
			Message: "cannot attend event in another domain",
			Status:  http.StatusForbidden,
		}
	}

	// Start a new transaction
	err = s.transactionManager.ExecuteInTransaction(func(tx *gorm.DB) error {
		// 5. Create EventAttendance record
		eventAttendance = entity.EventAttendance{
			FollowerID: follower.ID,
			EventID:    event.ID,
			AttendedAt: time.Now(),
		}
		err = s.eventAttendanceRepo.Create(ctx, &eventAttendance, tx)
		if err != nil {
			return err
		}

		follower.Points = follower.Points + event.PointsAwarded
		err = s.followerRepo.Update(ctx, &follower, tx) // Assuming Update uses the transaction
		if err != nil {
			return err
		}

		// 7. Create PointMutation record
		pointMutation := entity.PointMutation{
			FollowerID:  follower.ID,
			Amount:      event.PointsAwarded,
			SourceType:  entity.PointMutationSourceTypeEventAttendance,
			SourceID:    &eventAttendance.ID, // Link to the EventAttendance record
			Description: utility.PointOf(fmt.Sprintf("Points awarded for attending event: %s", event.Name)),
		}
		err = s.pointMutationRepo.Create(ctx, &pointMutation, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return eventAttendance, &errors.AppError{
			Message: "failed to attend event",
			Status:  http.StatusInternalServerError,
			Err:     err,
		}
	}

	return
}
