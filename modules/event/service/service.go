package service

import (
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/card"
	"github.com/PhantomX7/dhamma/modules/event"
	"github.com/PhantomX7/dhamma/modules/event_attendance"
	"github.com/PhantomX7/dhamma/modules/follower"
	"github.com/PhantomX7/dhamma/modules/point_mutation"
)

type service struct {
	eventRepo           event.Repository
	followerRepo        follower.Repository         // Add follower repository
	eventAttendanceRepo event_attendance.Repository // Add event_attendance repository
	pointMutationRepo   point_mutation.Repository   // Add point_mutation repository
	cardRepo            card.Repository             // Add card repository
	transactionManager  transaction_manager.Client
}

// New creates a new event service instance.
func New(
	eventRepo event.Repository,
	followerRepo follower.Repository,
	eventAttendanceRepo event_attendance.Repository,
	pointMutationRepo point_mutation.Repository,
	cardRepo card.Repository,
	transactionManager transaction_manager.Client,
) event.Service {
	return &service{
		eventRepo:           eventRepo,
		followerRepo:        followerRepo,
		eventAttendanceRepo: eventAttendanceRepo,
		pointMutationRepo:   pointMutationRepo,
		cardRepo:            cardRepo,
		transactionManager:  transactionManager,
	}
}
