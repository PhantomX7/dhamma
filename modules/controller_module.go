package modules

import (
	"go.uber.org/fx"

	authController "github.com/PhantomX7/dhamma/modules/auth/controller"
)

var ControllerModule = fx.Options(
	fx.Provide(
		authController.New,
	),
)
