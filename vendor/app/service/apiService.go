package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"app/model"
	"app/shared/email"
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
		var contentLength int
		var costTime int
		var statusCode int

		timeout := time.Duration(api.Timeout) * time.Millisecond
		client := http.Client{
			Timeout: timeout,
		}
		timeStart := time.Now()
		resp, err := client.Get(api.URL)

		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err)
			contentLength = -1
			costTime = -1
			if resp != nil {
				statusCode = resp.StatusCode
			} else {
				statusCode = -1
			}

		} else {
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)
			statusCode = resp.StatusCode
			contentLength = len(content)
			costTime = int(time.Since(timeStart) / time.Millisecond)
		}
		fmt.Printf("%s : %d\n", api.URL, statusCode)
		model.RequestCreate(api.ID, statusCode, costTime, contentLength)

		if api.FailMax == 1 && statusCode != 200 && api.AlertReceivers != "" {
			subject := "接口监控报警邮件"
			body := fmt.Sprintf("%s 已经连续请求失败已达预设(%d次)上限，请尽快验证服务", api.URL, api.FailMax)
			err = email.SendMail(api.AlertReceivers, subject, body)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
