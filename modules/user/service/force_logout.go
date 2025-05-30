package service

import (
	"context"
)

func (s *service) ForceLogout(ctx context.Context, userID uint64) (err error) {
	return s.refreshTokenRepo.InvalidateAllByUserID(ctx, userID)
}
