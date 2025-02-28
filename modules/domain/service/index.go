package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements domain.Service.
func (s *service) Index(pg *pagination.Pagination, ctx context.Context) (
	domains []entity.Domain, meta utility.PaginationMeta, err error,
) {
	domains, err = s.domainRepo.FindAll(pg, ctx)
	if err != nil {
		return
	}

	count, err := s.domainRepo.Count(pg, ctx)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
