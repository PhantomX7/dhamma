package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByCode(ctx context.Context, code string) (entity.Domain, error) {
	return r.base.FindOneByField(ctx, "code", code)
}
