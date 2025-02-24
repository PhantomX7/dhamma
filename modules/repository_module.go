package modules

import (
	"go.uber.org/fx"

	refreshTokenRepo "github.com/PhantomX7/dhamma/modules/refresh_token/repository"
	userRepo "github.com/PhantomX7/dhamma/modules/user/repository"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		refreshTokenRepo.New,
		userRepo.New,
	),
)
