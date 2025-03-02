package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
)

func (s *service) Update(roleID uint64, request request.RoleUpdateRequest, ctx context.Context) (role entity.Role, err error) {
	role, err = s.roleRepo.FindByID(roleID, ctx)
	if err != nil {
		return
	}

	err = copier.Copy(&role, &request)
	if err != nil {
		return
	}

	err = s.roleRepo.Update(&role, nil, ctx)
	if err != nil {
		return
	}

	return
}
