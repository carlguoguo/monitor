package service

import (
	"fmt"
	"net/http"

	"app/model"
	"app/shared/interval"
)

var apiJobMap = make(map[uint32]interval.Job)

// NewMonitor inits a job
func NewMonitor(api model.API) {
	job := interval.NewJob(requestGET(api), int(api.IntervalTime))
	apiJobMap[api.ID] = job
}

// StartMonitor creates a job and start monitor
func StartMonitor(api model.API) error {
	if err := model.APIUpdateStart(1, fmt.Sprint(api.ID)); err != nil {
		return err
	}
	// Re-init a new API object
	NewMonitor(api)
	job := apiJobMap[api.ID]
	job.Start()
	return nil
}

// PauseMonitor pause the monitor
func PauseMonitor(api model.API) error {
	if err := model.APIUpdateStart(0, fmt.Sprint(api.ID)); err != nil {
		return err
	}
	if job, ok := apiJobMap[api.ID]; ok {
		delete(apiJobMap, api.ID)
		job.Stop()
	}
	return nil
}

// StopMonitor stop the monitor
func StopMonitor(api model.API) error {
	job := apiJobMap[api.ID]
	delete(apiJobMap, api.ID)
	job.Stop()
	return nil
}

// RestartMonitor restarts the monitor
func RestartMonitor(api model.API) {
	PauseMonitor(api)
	StartMonitor(api)
}

func requestGET(api model.API) func() {
	return func() {
		resp, err := http.Get(api.URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		fmt.Printf("%s : %d\n", api.URL, resp.StatusCode)
	}
}
