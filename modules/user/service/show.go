package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements user.Service.
func (s *service) Show(userID uint64, ctx context.Context) (entity.User, error) {
	panic("unimplemented")
}
