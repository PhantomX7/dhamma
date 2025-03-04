package service

import (
	"context"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements user.Service.
func (s *service) Index(pg *pagination.Pagination, ctx context.Context) (
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

	users, err = s.userRepo.FindAll(pg, ctx)
	if err != nil {
		return
	}

	// only show specific domain if user have domain context
	if haveDomain {
		for i := range users {
			var domain entity.Domain
			domain, err = s.domainRepo.FindByID(domainID, ctx)
			if err != nil {
				return
			}

			users[i].Domains = []entity.Domain{domain}
		}
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
