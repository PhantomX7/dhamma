package modules

import (
	userService "github.com/PhantomX7/dhamma/modules/user/service"
	"go.uber.org/fx"

	authService "github.com/PhantomX7/dhamma/modules/auth/service"
	cronService "github.com/PhantomX7/dhamma/modules/cron/service"
	domainService "github.com/PhantomX7/dhamma/modules/domain/service"
	roleService "github.com/PhantomX7/dhamma/modules/role/service"
)

var ServiceModule = fx.Options(
	fx.Provide(
		authService.New,
		cronService.New,
		domainService.New,
		roleService.New,
		userService.New,
	),
)
