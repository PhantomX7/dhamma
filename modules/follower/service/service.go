package service

import (
	"github.com/PhantomX7/dhamma/modules/card" // Import card module
	"github.com/PhantomX7/dhamma/modules/follower"
)

type service struct {
	followerRepo follower.Repository
	cardRepo     card.Repository // Add card repository
}

// New creates a new follower service instance.
func New(
	followerRepo follower.Repository,
	cardRepo card.Repository, // Inject card repository
) follower.Service {
	return &service{
		followerRepo: followerRepo,
		cardRepo:     cardRepo, // Initialize card repository
	}
}
