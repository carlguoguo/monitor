package service

import (
	"fmt"
	"net/http"

	"app/model"
	"app/shared/interval"
)

var apiJobMap = make(map[uint32]interval.Job)

// StartMonitor creates a job and start monitor
func StartMonitor(api model.API) {
	job := interval.NewJob(RequestGET(api), int(api.IntervalTime))
	apiJobMap[api.ID] = job
	job.Start()
}

// StopMonitor stop the monitor
func StopMonitor(api model.API) {
	job := apiJobMap[api.ID]
	delete(apiJobMap, api.ID)
	job.Stop()
}

// RestartMonitor restarts the monitor
func RestartMonitor(api model.API) {
	StopMonitor(api)
	StartMonitor(api)
}

// RequestGET request a url
func RequestGET(api model.API) func() {
	return func() {
		resp, err := http.Get(api.URL)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		fmt.Printf("%s : %d\n", api.URL, resp.StatusCode)
	}
}
