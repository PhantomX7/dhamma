package controller

import (
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/health"
)

type controller struct {
	middleware *middleware.Middleware
}

func New(middleware *middleware.Middleware) health.Controller {
	return &controller{
		middleware: middleware,
	}
}