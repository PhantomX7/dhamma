package modules

import (
	"go.uber.org/fx"

	authService "github.com/PhantomX7/dhamma/modules/auth/service"
	cronService "github.com/PhantomX7/dhamma/modules/cron/service"
	domainService "github.com/PhantomX7/dhamma/modules/domain/service"
	permissionService "github.com/PhantomX7/dhamma/modules/permission/service"
	roleService "github.com/PhantomX7/dhamma/modules/role/service"
	userService "github.com/PhantomX7/dhamma/modules/user/service"
)

var ServiceModule = fx.Options(
	fx.Provide(
		authService.New,
		cronService.New,
		domainService.New,
		permissionService.New,
		roleService.New,
		userService.New,
	),
)
