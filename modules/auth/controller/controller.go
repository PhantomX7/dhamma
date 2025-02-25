package controller

import (
	"github.com/PhantomX7/dhamma/modules/auth"
)

type controller struct {
	authService auth.Service
}

func New(authService auth.Service) auth.Controller {
	return &controller{
		authService: authService,
	}
}
