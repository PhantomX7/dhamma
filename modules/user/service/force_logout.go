package service

import (
	"context"
)

func (s *service) ForceLogout(userID uint64, ctx context.Context) (err error) {
	return s.refreshTokenRepo.InvalidateAllByUserID(userID, ctx)
}
