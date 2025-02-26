package modules

import (
	"go.uber.org/fx"

	domainRepo "github.com/PhantomX7/dhamma/modules/domain/repository"
	refreshTokenRepo "github.com/PhantomX7/dhamma/modules/refresh_token/repository"
	roleRepo "github.com/PhantomX7/dhamma/modules/role/repository"
	userRepo "github.com/PhantomX7/dhamma/modules/user/repository"
	userRoleRepo "github.com/PhantomX7/dhamma/modules/user_role/repository"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		domainRepo.New,
		refreshTokenRepo.New,
		roleRepo.New,
		userRepo.New,
		userRoleRepo.New,
	),
)
