package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements domain.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	domains []entity.Domain, meta utility.PaginationMeta, err error,
) {
	domains, err = s.domainRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.domainRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
