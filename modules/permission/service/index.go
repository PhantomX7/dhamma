package service

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements permission.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	permissions []entity.Permission, meta utility.PaginationMeta, err error,
) {
	//pg.AddCustomScope(func(db *gorm.DB) *gorm.DB {
	//	return db.Where("domain_id = ?", *contextValues.DomainID)
	//})
	permissions, err = s.permissionRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.permissionRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
