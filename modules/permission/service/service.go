package service

import (
	"github.com/PhantomX7/dhamma/modules/permission"
)

type service struct {
	permissionRepo permission.Repository
}

func New(
	permissionRepo permission.Repository,
) permission.Service {
	return &service{
		permissionRepo: permissionRepo,
	}
}
