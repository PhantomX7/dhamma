package service

import (
	"context"
)

func (u *service) ClearRefreshToken() (err error) {
	err = u.refreshTokenRepo.DeleteInvalidToken(context.Background())
	return
}
