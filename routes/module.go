package routes

import (
	"github.com/PhantomX7/dhamma/routes/admin"
	"github.com/PhantomX7/dhamma/routes/domain"

	"go.uber.org/fx"
)

var Module = fx.Invoke(
	admin.AuthRoute,
	admin.DomainRoute,
	admin.EventRoute,
	admin.FollowerRoute,
	admin.PermissionRoute,
	admin.UserRoute,
	admin.RoleRoute,

	// domain specific route
	domain.AuthRoute,
	domain.EventRoute,
	domain.FollowerRoute,
	domain.PermissionRoute,
	domain.UserRoute,
	domain.RoleRoute,

	Universal,
)
