package service

import (
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/libs/gocache"
	"github.com/PhantomX7/dhamma/modules/role"
)

type service struct {
	roleRepo role.Repository
	casbin   casbin.Client
	cache    gocache.Client
}

func New(
	roleRepo role.Repository,
	casbin casbin.Client,
	cache gocache.Client,
) role.Service {
	return &service{
		roleRepo: roleRepo,
		casbin:   casbin,
		cache:    cache,
	}
}
