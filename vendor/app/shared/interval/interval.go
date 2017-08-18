package interval

import (
	"time"
)

// Job is job what will intervally do
type Job struct {
	JobFunc      func()
	IntervalTime *time.Ticker
	Quit         chan bool
}

func (job *Job) run() {
	go func() {
		for {
			select {
			case <-job.IntervalTime.C:
				job.JobFunc()
			case <-job.Quit:
				job.IntervalTime.Stop()
				return
			}
		}
	}()
}

// NewJob init a job
func NewJob(jobFunc func(), intervalTime int) Job {
	job := Job{
		JobFunc:      jobFunc,
		IntervalTime: time.NewTicker(time.Duration(intervalTime) * time.Second),
		Quit:         make(chan bool, 1),
	}
	return job
}

// Start starts a job
func (job *Job) Start() {
	job.run()
}

// Stop stops a job
func (job *Job) Stop() {
	job.Quit <- true
}
