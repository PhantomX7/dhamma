package cron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Service interface {
	ClearRefreshToken() error
}

func New() *gocron.Scheduler {

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Sprint("error starting cron: ", err.Error()))
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
			},
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	return &s
}
