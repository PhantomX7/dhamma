package modules

import (
	"go.uber.org/fx"

	authService "github.com/PhantomX7/dhamma/modules/auth/service"
)

var ServiceModule = fx.Options(
	fx.Provide(
		authService.New,
	),
)
