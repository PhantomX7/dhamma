package controller

import (
	"github.com/PhantomX7/dhamma/modules/follower"
)

type controller struct {
	followerService follower.Service
}

func New(followerService follower.Service) follower.Controller {
	return &controller{
		followerService: followerService,
	}
}
