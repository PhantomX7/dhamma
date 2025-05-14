package modules

import (
	"go.uber.org/fx"

	cardRepo "github.com/PhantomX7/dhamma/modules/card/repository"
	domainRepo "github.com/PhantomX7/dhamma/modules/domain/repository"
	followerRepo "github.com/PhantomX7/dhamma/modules/follower/repository"
	permissionRepo "github.com/PhantomX7/dhamma/modules/permission/repository"
	refreshTokenRepo "github.com/PhantomX7/dhamma/modules/refresh_token/repository"
	roleRepo "github.com/PhantomX7/dhamma/modules/role/repository"
	userRepo "github.com/PhantomX7/dhamma/modules/user/repository"
	userDomainRepo "github.com/PhantomX7/dhamma/modules/user_domain/repository"
	userRoleRepo "github.com/PhantomX7/dhamma/modules/user_role/repository"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		cardRepo.New,
		domainRepo.New,
		followerRepo.New,
		permissionRepo.New,
		refreshTokenRepo.New,
		roleRepo.New,
		userRepo.New,
		userDomainRepo.New,
		userRoleRepo.New,
	),
)
