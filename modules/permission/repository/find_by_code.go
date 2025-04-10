package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByCode(ctx context.Context, permissionCode string) (entity.Permission, error) {
	return r.base.FindOneByField(ctx, "code", permissionCode)
}
