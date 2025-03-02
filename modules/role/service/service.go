package service

import (
	"github.com/PhantomX7/dhamma/modules/role"
)

type service struct {
	roleRepo role.Repository
}

func New(
	roleRepo role.Repository,
) role.Service {
	return &service{
		roleRepo: roleRepo,
	}
}
