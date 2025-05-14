package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements follower.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	followers []entity.Follower, meta utility.PaginationMeta, err error,
) {

	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	// Combine all scopes into a single AddCustomScope call
	pg.AddCustomScope(
		func(db *gorm.DB) *gorm.DB {
			// Join users -> user_domains -> domains. Alias domains as 'd' for the filter.
			// Use Group("users.id") to prevent duplicate user rows if a user belongs to multiple domains
			// when only filtering by domain name (not context domain ID).
			return db.
				Joins("LEFT JOIN cards Card ON Card.follower_id = followers.id").
				Preload("Cards").
				Group("followers.id") // Group by user ID to avoid duplicates from joins
		},
		// Base join and preload
		func(db *gorm.DB) *gorm.DB {
			if contextValues.DomainID != nil {
				return db.Where("domain_id = ?", *contextValues.DomainID)
			}
			return db.
				Joins("Domain")
		},
	)

	followers, err = s.followerRepo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.followerRepo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
