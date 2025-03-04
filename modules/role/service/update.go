package service

import (
	"context"
	"errors"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
)

func (s *service) Update(ctx context.Context, roleID uint64, request request.RoleUpdateRequest) (role entity.Role, err error) {
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
			return entity.Role{}, errors.New("you are not allowed to create role for another domain")
		}
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
