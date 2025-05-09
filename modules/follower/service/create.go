package service

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/follower/dto/request"
	"github.com/PhantomX7/dhamma/utility" // Import the main utility package
)

// Create handles the logic for creating a new follower.
// It includes a domain context check to ensure the follower is created for the allowed domain.
func (s *service) Create(ctx context.Context, request request.FollowerCreateRequest) (follower entity.Follower, err error) {
	// Perform domain context check using DomainID from the request and the generic helper
	_, err = utility.CheckDomainContext(ctx, request.DomainID, "follower", "create")
	if err != nil {
		return
	}

	follower = entity.Follower{
		Points:   0, // Default points
		DomainID: request.DomainID,
	}

	err = copier.Copy(&follower, &request)
	if err != nil {
		return
	}

	err = s.followerRepo.Create(ctx, &follower, nil)
	if err != nil {
		return
	}

	return
}
