package service

import (
	"context"
	"fmt"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/permission/dto/request"
	"github.com/jinzhu/copier"
)

func (s *service) Create(ctx context.Context, request request.PermissionCreateRequest) (permission entity.Permission, err error) {
	err = copier.Copy(&permission, &request)
	if err != nil {
		return
	}

	permission.Code = fmt.Sprintf("web:%s:%s", permission.Object, permission.Action)
	permission.Type = constants.EnumPermissionTypeWeb

	err = s.permissionRepo.Create(ctx, &permission, nil)
	if err != nil {
		return
	}

	return
}
