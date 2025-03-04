package service

import (
	"context"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements role.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	roles []entity.Role, meta utility.PaginationMeta, err error,
) {
	haveDomain, domainID := utility.GetDomainIDFromContext(ctx)

	// only query specific domain
	if haveDomain {
		pg.AddCustomScope(func(db *gorm.DB) *gorm.DB {
			return db.Where("domain_id = ?", domainID)
		})
	}

	roles, err = s.roleRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.roleRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
