package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByNameAndDomainID(
	ctx context.Context, name string, domainID uint64,
) (entity.Role, error) {
	return r.base.FindOneByFields(ctx, map[string]any{"name": name, "domain_id": domainID}, "Domain")
}
