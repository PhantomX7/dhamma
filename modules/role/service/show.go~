package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements role.Service
func (s *service) Show(roleID uint64, ctx context.Context) (role entity.Role, err error) {
	role, err = s.roleRepo.FindByID(roleID, ctx)
	if err != nil {
		return
	}

	return
}
