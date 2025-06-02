package modules

import (
	"go.uber.org/fx"

	authController "github.com/PhantomX7/dhamma/modules/auth/controller"
	chatTemplateController "github.com/PhantomX7/dhamma/modules/chat_template/controller"
	cronController "github.com/PhantomX7/dhamma/modules/cron/controller"
	domainController "github.com/PhantomX7/dhamma/modules/domain/controller"
	eventController "github.com/PhantomX7/dhamma/modules/event/controller"
	eventAttendanceController "github.com/PhantomX7/dhamma/modules/event_attendance/controller"
	followerController "github.com/PhantomX7/dhamma/modules/follower/controller"
	healthController "github.com/PhantomX7/dhamma/modules/health/controller"
	permissionController "github.com/PhantomX7/dhamma/modules/permission/controller"
	pointMutationController "github.com/PhantomX7/dhamma/modules/point_mutation/controller"
	roleController "github.com/PhantomX7/dhamma/modules/role/controller"
	userController "github.com/PhantomX7/dhamma/modules/user/controller"
)

var ControllerModule = fx.Options(
	fx.Provide(
		authController.New,
		chatTemplateController.New,
		cronController.New,
		domainController.New,
		eventController.New,
		eventAttendanceController.New,
		followerController.New,
		healthController.New,
		permissionController.New,
		pointMutationController.New,
		roleController.New,
		userController.New,
	),
)
