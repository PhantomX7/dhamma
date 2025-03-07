package service

import (
	"context"
	"github.com/PhantomX7/dhamma/entity"
)

// Show implements permission.Service
func (s *service) Show(ctx context.Context, permissionID uint64) (permission entity.Permission, err error) {
	permission, err = s.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return
	}

	return
}
