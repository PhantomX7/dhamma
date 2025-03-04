package cron

type Service interface {
	ClearRefreshToken() error
}
