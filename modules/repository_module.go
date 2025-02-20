package modules

import (
	"go.uber.org/fx"

	userRepo "github.com/PhantomX7/dhamma/modules/user/repository"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		userRepo.New,
	),
)
