package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

func (r *repository) FindByUserIDAndDomainID(ctx context.Context, userID uint64, domainID uint64, preloads ...string) ([]entity.UserRole, error) {
	return r.base.FindByFields(ctx, map[string]any{"user_id": userID, "domain_id": domainID}, preloads...)
}
