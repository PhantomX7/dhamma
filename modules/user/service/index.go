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
		// Base join and preload - Ensure necessary tables are joined for filtering and preloading
		func(db *gorm.DB) *gorm.DB {
			// Join users -> user_domains -> domains. Alias domains as 'd' for the filter.
			// Use Group("users.id") to prevent duplicate user rows if a user belongs to multiple domains
			// when only filtering by domain name (not context domain ID).
			return db.
				Joins("LEFT JOIN user_domains ud ON ud.user_id = users.id").
				Joins("LEFT JOIN domains d ON ud.domain_id = d.id").
				Preload("Domains"). // Preload associated domains
				Group("users.id")   // Group by user ID to avoid duplicates from joins
		},
		// Domain filter scope (based on context)
		func(db *gorm.DB) *gorm.DB {
			// Apply context domain filter if present. The join is already done above.
			if contextValues.DomainID != nil {
				// Use the alias 'ud' defined in the base join
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

	// Find users using the repository method which applies pagination filters
	users, err = s.userRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	// Post-processing: only show specific domain if user has domain context
	// This overrides the preloaded domains if a specific context domain is set.
	if contextValues.DomainID != nil {
		var domain entity.Domain
		domain, err = s.domainRepo.FindByID(ctx, *contextValues.DomainID)
		if err != nil {
			// Handle error, maybe return or log
			return
		}

		// Ensure the Domains slice exists and set the specific domain
		for i := range users {
			users[i].Domains = []entity.Domain{domain}
		}
	}

	// Count total matching users using the repository method with the same filters
	count, err := s.userRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	// Set pagination metadata
	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
