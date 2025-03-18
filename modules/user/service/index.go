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
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	// Combine all scopes into a single AddCustomScope call
	pg.AddCustomScope(
		// Base join and preload
		func(db *gorm.DB) *gorm.DB {
			return db.Joins("LEFT JOIN user_domains ud ON ud.user_id = users.id").
				Preload("Domains")
		},
		// Domain filter scope
		func(db *gorm.DB) *gorm.DB {
			if contextValues.DomainID != nil {
				return db.Where("ud.domain_id = ?", *contextValues.DomainID)
			}
			return db
		},
		// Super admin filter scope
		func(db *gorm.DB) *gorm.DB {
			if !contextValues.IsRoot {
				return db.Where("users.is_super_admin = ?", false)
			}
			return db
		},
	)

	users, err = s.userRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	// only show specific domain if user have domain context
	if contextValues.DomainID != nil {
		var domain entity.Domain
		domain, err = s.domainRepo.FindByID(ctx, *contextValues.DomainID)
		if err != nil {
			return
		}

		for i := range users {
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
