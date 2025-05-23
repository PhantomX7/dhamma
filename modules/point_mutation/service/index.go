package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements point_mutation.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	pointMutations []entity.PointMutation, meta utility.PaginationMeta, err error,
) {
	pointMutations, err = s.pointMutationRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.pointMutationRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
