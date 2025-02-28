package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements user.Service.
func (s *service) Index(pg *pagination.Pagination, ctx context.Context) (
	users []entity.User, meta utility.PaginationMeta, err error,
) {
	users, err = s.userRepo.FindAll(pg, ctx)
	if err != nil {
		return
	}

	count, err := s.userRepo.Count(pg, ctx)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
