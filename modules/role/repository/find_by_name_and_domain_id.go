package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

func (r *repository) FindByNameAndDomainID(
	ctx context.Context, name string, domainID uint64,
) (entity.Role, error) {
	db := r.db.WithContext(ctx).Preload("Domain")

	var role entity.Role
	result := db.Where("name = ? AND domain_id = ?", name, domainID).Take(&role)

	if result.Error != nil {
		return role, utility.WrapError(utility.ErrNotFound, "role with name %s and domain ID %d not found", name, domainID)
	}

	return role, nil
}
