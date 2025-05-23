package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements point_mutation.Service
func (s *service) Show(ctx context.Context, pointMutationID uint64) (pointMutation entity.PointMutation, err error) {
	pointMutation, err = s.pointMutationRepo.FindByID(ctx, pointMutationID)
	if err != nil {
		return
	}

	return
}
