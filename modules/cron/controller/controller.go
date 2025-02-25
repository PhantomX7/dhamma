package controller

import (
	"fmt"
	"time"

	"github.com/PhantomX7/dhamma/modules/cron"

	"github.com/go-co-op/gocron/v2"
)

func New(cronService cron.Service) gocron.Scheduler {

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Sprint("error starting cron: ", err.Error()))
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(2*time.Second),
		gocron.NewTask(cronService.ClearRefreshToken),
	)
	if err != nil {
		panic(fmt.Sprint("error initializing job clear refresh token: ", err.Error()))
	}
	// each job has a unique id

	return s
}
