package service

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/constants/permissions"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/permission/dto/request"
)

func (s *service) Create(ctx context.Context, request request.PermissionCreateRequest) (permission entity.Permission, err error) {
	err = copier.Copy(&permission, &request)
	if err != nil {
		return
	}

	permission.Code = fmt.Sprintf("web:%s/%s", permission.Object, permission.Action)
	permission.Type = permissions.PermissionTypeWeb

	err = s.permissionRepo.Create(ctx, &permission, nil)
	if err != nil {
		return
	}

	return
}
