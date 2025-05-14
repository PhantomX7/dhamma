package service

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) Update(ctx context.Context, roleID uint64, request request.RoleUpdateRequest) (role entity.Role, err error) {
	role, err = s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	// Perform domain context check using DomainID from the request and the generic helper
	_, err = utility.CheckDomainContext(ctx, role.DomainID, "role", "update")
	if err != nil {
		return
	}

	err = copier.Copy(&role, &request)
	if err != nil {
		return
	}

	err = s.roleRepo.Update(ctx, &role, nil)
	if err != nil {
		return
	}

	return
}
