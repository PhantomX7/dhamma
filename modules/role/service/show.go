package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
)

// Show implements role.Service
func (s *service) Show(ctx context.Context, roleID uint64) (role entity.Role, err error) {
	role, err = s.roleRepo.FindByID(ctx, roleID, "Domain")
	if err != nil {
		return
	}

	// Perform domain context check using DomainID from the request and the generic helper
	_, err = utility.CheckDomainContext(ctx, role.DomainID, "role", "show")
	if err != nil {
		return
	}

	role.Permissions = s.casbin.GetRolePermissions(role.ID, role.DomainID)

	return
}
