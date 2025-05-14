package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/follower/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

// AddCard handles the logic for adding a new card to an existing follower.
func (s *service) AddCard(ctx context.Context, followerID uint64, req request.FollowerAddCardRequest) (card entity.Card, err error) {
	// Find the follower
	follower, err := s.followerRepo.FindByID(ctx, followerID)
	if err != nil {
		return
	}

	// Perform domain context check for the follower
	_, err = utility.CheckDomainContext(ctx, follower.DomainID, "follower", "add card to")
	if err != nil {
		return
	}

	// Prepare the new card entity
	card = entity.Card{
		FollowerID: followerID,
		DomainID:   follower.DomainID, // Associate card with the follower's domain
		Code:       req.Code,
	}

	// Create the card
	err = s.cardRepo.Create(ctx, &card, nil)
	if err != nil {
		return
	}

	return
}
