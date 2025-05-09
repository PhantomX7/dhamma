package service

import (
	"github.com/PhantomX7/dhamma/modules/follower"
)

type service struct {
	followerRepo follower.Repository
}

func New(
	followerRepo follower.Repository,
) follower.Service {
	return &service{
		followerRepo: followerRepo,
	}
}
