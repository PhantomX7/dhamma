package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements role.Service.
func (s *service) Index(pg *pagination.Pagination, ctx context.Context) (
	roles []entity.Role, meta utility.PaginationMeta, err error,
) {
	roles, err = s.roleRepo.FindAll(pg, ctx)
	if err != nil {
		return
	}

	count, err := s.roleRepo.Count(pg, ctx)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
