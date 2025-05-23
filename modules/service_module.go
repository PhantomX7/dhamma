package modules

import (
	"go.uber.org/fx"

	authService "github.com/PhantomX7/dhamma/modules/auth/service"
	cronService "github.com/PhantomX7/dhamma/modules/cron/service"
	domainService "github.com/PhantomX7/dhamma/modules/domain/service"
	eventService "github.com/PhantomX7/dhamma/modules/event/service"
	eventAttendanceService "github.com/PhantomX7/dhamma/modules/event_attendance/service"
	followerService "github.com/PhantomX7/dhamma/modules/follower/service"
	permissionService "github.com/PhantomX7/dhamma/modules/permission/service"
	pointMutationService "github.com/PhantomX7/dhamma/modules/point_mutation/service"
	roleService "github.com/PhantomX7/dhamma/modules/role/service"
	userService "github.com/PhantomX7/dhamma/modules/user/service"
)

var ServiceModule = fx.Options(
	fx.Provide(
		authService.New,
		cronService.New,
		domainService.New,
		eventService.New,
		eventAttendanceService.New,
		followerService.New,
		permissionService.New,
		pointMutationService.New,
		roleService.New,
		userService.New,
	),
)
