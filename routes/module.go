package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Invoke(
	Auth,
	Domain,
	Universal,
	User,
)
