package routes

import (
	"github.com/PhantomX7/dhamma/routes/admin"
	"github.com/PhantomX7/dhamma/routes/domain"

	"go.uber.org/fx"
)

var Module = fx.Invoke(
	admin.Auth,
	admin.Domain,
	admin.User,
	admin.Role,

	// domain specific route
	domain.Auth,
	domain.User,

	Universal,
)
