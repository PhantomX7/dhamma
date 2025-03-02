package controller

import (
	"github.com/PhantomX7/dhamma/modules/role"
)

type controller struct {
	roleService role.Service
}

func New(roleService role.Service) role.Controller {
	return &controller{
		roleService: roleService,
	}
}
