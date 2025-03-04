package service

import (
	"context"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
)

func (s *service) Update(ctx context.Context, roleID uint64, request request.RoleUpdateRequest) (role entity.Role, err error) {
	hasDomain, domainID := utility.GetDomainIDFromContext(ctx)

	role, err = s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return
	}

	if hasDomain {
		if role.DomainID != domainID {
			return entity.Role{}, utility.LogError("you are not allowed to create role for another domain", nil)
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
