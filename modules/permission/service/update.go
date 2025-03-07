package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/permission/dto/request"
)

func (s *service) Update(ctx context.Context, permissionID uint64, request request.PermissionUpdateRequest) (permission entity.Permission, err error) {
	permission, err = s.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return
	}

	err = copier.Copy(&permission, &request)
	if err != nil {
		return
	}

	err = s.permissionRepo.Update(ctx, &permission, nil)
	if err != nil {
		return
	}

	return
}
