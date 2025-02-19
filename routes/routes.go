package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Invoke(
	User,
	Universal,
)
