package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var q = make(chan Job, 1)

func main() {
	defer close(q)

	// workers
	for id := range 10 {
		go func(id int) {
			for j := range q {
				j.Run(id)
			}
		}(id + 1)
	}

	c := cron.New(cron.WithSeconds())
	// queue
	for i := range 10 {
		c.AddJob("@every 1s", NewJob(fmt.Sprintf("job-%02d", i)))
	}
	c.Run()
}

// job wrapper
type JobWrapper Job

func (j JobWrapper) Run() {
	q <- Job(j)
}

func WrapJob(j Job) JobWrapper {
	return JobWrapper(j)
}

// actual job
type Job string

func NewJob(name string) JobWrapper {
	return WrapJob(Job(name))
}

func (j Job) Run(workerID int) {
	fmt.Println(workerID, j, time.Now().UnixNano())
}
