package routes

import (
	"github.com/PhantomX7/dhamma/routes/admin"
	"github.com/PhantomX7/dhamma/routes/domain"

	"go.uber.org/fx"
)

var Module = fx.Invoke(
	admin.AuthRoute,
	admin.ChatTemplateRoute,
	admin.DomainRoute,
	admin.EventAttendanceRoute,
	admin.EventRoute,
	admin.FollowerRoute,
	admin.HealthRoute,
	admin.PermissionRoute,
	admin.PointMutationRoute,
	admin.UserRoute,
	admin.RoleRoute,

	// domain specific route
	domain.AuthRoute,
	domain.ChatTemplateRoute,
	domain.EventAttendanceRoute,
	domain.EventRoute,
	domain.FollowerRoute,
	domain.PermissionRoute,
	domain.PointMutationRoute,
	domain.UserRoute,
	domain.RoleRoute,

	Universal,
)
