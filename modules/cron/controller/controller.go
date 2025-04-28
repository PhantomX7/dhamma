package controller

import (
	"time"

	"github.com/PhantomX7/dhamma/modules/cron"
	"github.com/PhantomX7/dhamma/utility/logger"
	"go.uber.org/zap"

	"github.com/go-co-op/gocron/v2"
)

func New(cronService cron.Service) gocron.Scheduler {

	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Get().Panic("error creating cron scheduler", zap.Error(err))
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(2*time.Second),
		gocron.NewTask(cronService.ClearRefreshToken),
	)
	if err != nil {
		logger.Get().Panic("error creating cron job for clear refresh token", zap.Error(err))
	}
	// each job has a unique id

	return s
}
