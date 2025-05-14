package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility" // Import the main utility package
)

// Show handles the logic for retrieving a specific follower.
// It includes a domain context check to ensure the follower belongs to the allowed domain.
func (s *service) Show(ctx context.Context, followerID uint64) (follower entity.Follower, err error) {
	follower, err = s.followerRepo.FindByID(ctx, followerID, "Domain", "Cards")
	if err != nil {
		return
	}

	// Perform domain context check using the generic helper
	_, err = utility.CheckDomainContext(ctx, follower.DomainID, "follower", "get")
	if err != nil {
		return
	}

	return
}
