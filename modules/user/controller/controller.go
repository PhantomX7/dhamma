package controller

import (
	"github.com/PhantomX7/dhamma/modules/user"
)

type controller struct {
	userService user.Service
}

func New(userService user.Service) user.Controller {
	return &controller{
		userService: userService,
	}
}
