package service

import (
	"context"
	// net/http is no longer needed here if AppError is handled by the utility function
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/follower/dto/request"
	"github.com/PhantomX7/dhamma/utility" // Import the main utility package
	// "github.com/PhantomX7/dhamma/utility/errors" // No longer needed directly here for AppError
)

// Update handles the logic for updating an existing follower.
// It includes a domain context check to ensure the follower belongs to the allowed domain.
func (s *service) Update(ctx context.Context, followerID uint64, request request.FollowerUpdateRequest) (follower entity.Follower, err error) {
	follower, err = s.followerRepo.FindByID(ctx, followerID)
	if err != nil {
		return
	}

	// Perform domain context check using the generic helper
	_, err = utility.CheckDomainContext(ctx, follower.DomainID, "follower", "update")
	if err != nil {
		return
	}

	err = copier.Copy(&follower, &request)
	if err != nil {
		return
	}

	err = s.followerRepo.Update(ctx, &follower, nil)
	if err != nil {
		return
	}

	return
}
