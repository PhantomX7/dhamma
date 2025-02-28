package modules

import (
	"go.uber.org/fx"

	authController "github.com/PhantomX7/dhamma/modules/auth/controller"
	cronController "github.com/PhantomX7/dhamma/modules/cron/controller"
	domainController "github.com/PhantomX7/dhamma/modules/domain/controller"
	userController "github.com/PhantomX7/dhamma/modules/user/controller"
)

var ControllerModule = fx.Options(
	fx.Provide(
		authController.New,
		cronController.New,
		domainController.New,
		userController.New,
	),
)
