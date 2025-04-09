package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) GetDomainRoles(ctx context.Context, domainID uint64) (entity.Domain, error) {
	// TODO: fix this approach
	return r.base.FindByID(ctx, domainID)
}
