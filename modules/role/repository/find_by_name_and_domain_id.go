package repository

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByNameAndDomainID(
	name string, domainID uint64, ctx context.Context,
) (roleM entity.Role, err error) {

	err = r.db.WithContext(ctx).
		Preload("Domain").
		Where("name = ? AND domain_id = ?", name, domainID).
		Take(&roleM).Error
	if err != nil {
		err = utility.LogError("error find role by name and domain id", err)
		return
	}

	return
}
