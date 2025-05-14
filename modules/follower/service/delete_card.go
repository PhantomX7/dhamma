package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// DeleteCard handles the logic for deleting a card associated with a follower.
func (s *service) DeleteCard(ctx context.Context, followerID uint64, cardID uint64) (err error) {
	// Find the follower
	follower, err := s.followerRepo.FindByID(ctx, followerID)
	if err != nil {
		return
	}

	// Perform domain context check for the follower
	_, err = utility.CheckDomainContext(ctx, follower.DomainID, "follower", "delete card from")
	if err != nil {
		return
	}

	// Find the card
	card, err := s.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return
	}

	// Validate card ownership and domain
	if card.FollowerID != followerID {
		return &errors.AppError{
			Message: fmt.Sprintf("card with id %d does not belong to follower %d", cardID, followerID),
			Status:  http.StatusForbidden,
		}
	}
	// Double check domain consistency, though follower check should cover it.
	if card.DomainID != follower.DomainID {
		return &errors.AppError{
			Message: fmt.Sprintf("card with id %d is not in the same domain as follower %d", cardID, followerID),
			Status:  http.StatusForbidden,
		}
	}

	// Delete the card
	err = s.cardRepo.Delete(ctx, &card, nil)
	if err != nil {
		return
	}

	return
}
