package service

import (
	"context"
	"errors"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (s *service) AddPermissions(ctx context.Context, roleID uint64, request request.RoleAddPermissionsRequest) (role entity.Role, err error) {
	// Get value from context
	contextValues, err := utility.ValuesFromContext(ctx)
	if err != nil {
		return
	}

	role, err = s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	if contextValues.DomainID != nil {
		if role.DomainID != *contextValues.DomainID {
			return entity.Role{}, errors.New("you are not allowed to add permissions to role for another domain")
		}
	}

	s.casbin.AddRolePermissions(role.ID, role.DomainID, request.Permissions)

	return
}
