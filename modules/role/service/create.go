package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/role/dto/request"
)

func (s *service) Create(request request.RoleCreateRequest, ctx context.Context) (role entity.Role, err error) {
	role = entity.Role{
		IsActive: true,
	}

	err = copier.Copy(&role, &request)
	if err != nil {
		return
	}

	err = s.roleRepo.Create(&role, nil, ctx)
	if err != nil {
		return
	}

	return
}
