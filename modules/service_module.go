package modules

import (
	"go.uber.org/fx"

	authService "github.com/PhantomX7/dhamma/modules/auth/service"
	userService "github.com/PhantomX7/dhamma/modules/user/service"
)

var ServiceModule = fx.Options(
	fx.Provide(
		authService.New,
		userService.New,
	),
)
