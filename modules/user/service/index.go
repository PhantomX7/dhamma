package service

import (
	"context"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements user.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	users []entity.User, meta utility.PaginationMeta, err error,
) {
	haveDomain, domainID := utility.GetDomainIDFromContext(ctx)

	pg.AddCustomScope(func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN user_domains ud ON ud.user_id = users.id").
			Preload("Domains")
	})

	// only query specific domain
	if haveDomain {
		pg.AddCustomScope(func(db *gorm.DB) *gorm.DB {
			return db.Where("ud.domain_id = ?", domainID)
		})
	}

	users, err = s.userRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	// only show specific domain if user have domain context
	if haveDomain {
		for i := range users {
			var domain entity.Domain
			domain, err = s.domainRepo.FindByID(ctx, domainID)
			if err != nil {
				return
			}

			users[i].Domains = []entity.Domain{domain}
		}
	}

	count, err := s.userRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
