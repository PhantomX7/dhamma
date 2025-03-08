package controller

import (
	"github.com/PhantomX7/dhamma/modules/permission"
)

type controller struct {
	permissionService permission.Service
}

func New(permissionService permission.Service) permission.Controller {
	return &controller{
		permissionService: permissionService,
	}
}
