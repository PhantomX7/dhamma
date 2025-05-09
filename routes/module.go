package routes

import (
	"github.com/PhantomX7/dhamma/routes/admin"
	"github.com/PhantomX7/dhamma/routes/domain"

	"go.uber.org/fx"
)

var Module = fx.Invoke(
	admin.AuthRoute,
	admin.DomainRoute,
	admin.UserRoute,
	admin.RoleRoute,
	admin.PermissionRoute,
	admin.FollowerRoute,

	// domain specific route
	domain.AuthRoute,
	domain.UserRoute,
	domain.RoleRoute,
	domain.PermissionRoute,

	Universal,
)
