package controller

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	"app/model"
	"app/service"
	"app/shared/session"
	"app/shared/view"
)

// APICreatePOST handles the api creation form submission
func APICreatePOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	url := r.FormValue("url")
	alias := r.FormValue("alias")
	intervalTime, _ := strconv.Atoi(r.FormValue("interval_time"))
	alertReceivers := r.FormValue("receivers")
	timeout, _ := strconv.Atoi(r.FormValue("timeout"))
	failMax, _ := strconv.Atoi(r.FormValue("fail_threshold"))
	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get database result
	newAPI, creationErr := model.APICreateAndGet(url, intervalTime, userID, alias, alertReceivers, timeout, failMax)
	if creationErr != nil {
		log.Println(creationErr)
		sess.AddFlash(view.Flash{Message: "服务器错误，请重试!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{Message: "新接口已经添加!", Class: view.FlashSuccess})
		sess.Save(r, w)
		model.APIStatusCreate(fmt.Sprintf("%d", newAPI.ID))
		service.NewMonitor(newAPI)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Direct to index
	IndexGET(w, r)
}

// APIUpdatePost handles the api update form submission
func APIUpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get form values
	url := r.FormValue("url")
	alias := r.FormValue("alias")
	alertReceivers := r.FormValue("receivers")
	timeout, _ := strconv.Atoi(r.FormValue("timeout"))
	failMax, _ := strconv.Atoi(r.FormValue("fail_threshold"))
	intervalTime, _ := strconv.Atoi(r.FormValue("interval_time"))

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")

	// Get database result
	newAPI, updateErr := model.APIUpdateAndReturn(url, intervalTime, userID, alias, alertReceivers, timeout, failMax, apiID)
	// Will only error if there is a problem with the query
	if updateErr != nil {
		log.Println(updateErr)
		sess.AddFlash(view.Flash{Message: "服务器错误，请重试!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{Message: "接口修改成功!", Class: view.FlashSuccess})
		sess.Save(r, w)
		service.PauseMonitor(newAPI)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Direct to index
	IndexGET(w, r)
}

// APIDeleteGet handles the api deletion
func APIDeleteGet(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")

	// Get database result
	deletedAPI, err := model.APIDeleteAndReturn(apiID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{Message: "服务器错误，请重试!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{Message: "接口删除成功!", Class: view.FlashSuccess})
		service.StopMonitor(deletedAPI)
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

// APIDetailGet display the api detail
func APIDetailGet(w http.ResponseWriter, r *http.Request) {
	session := session.Instance(r)
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")
	api, apiErr := model.APIByID(apiID)
	apiStatus, apiStatusErr := model.APIStatusByID(apiID)
	requests, requestErr := model.RequestByAPIID(apiID, 10, true)
	if apiErr != nil || apiStatusErr != nil || requestErr != nil {
		log.Println(apiErr)
		log.Println(apiStatusErr)
		log.Println(requestErr)
		session.AddFlash(view.Flash{Message: "服务器错误，请重试!", Class: view.FlashError})
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	v := view.New(r)
	v.Name = "api/detail"
	v.Vars["username"] = session.Values["username"]
	v.Vars["api"] = api
	v.Vars["apiStatus"] = apiStatus
	v.Vars["requests"] = requests
	v.Render(w)
}

type point struct {
	Date  string
	Value int
}

func (p *point) marshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`["`)
	buf.WriteString(p.Date)
	buf.WriteString(`",`)
	buf.WriteString(strconv.Itoa(p.Value))
	buf.WriteRune(']')

	return buf.Bytes(), nil
}

// APIRequestDetail return api request detail data
func APIRequestDetail(w http.ResponseWriter, r *http.Request) {
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")
	requests, requestErr := model.RequestByAPIID(apiID, -1, false)
	if requestErr != nil {
		http.Error(w, requestErr.Error(), http.StatusInternalServerError)
		return
	}
	result := []byte{}
	total := len(requests)
	for index, request := range requests {
		p := point{Date: request.RequestTime.Format("2006-01-02 15:04:05"), Value: request.Cost}
		if index == 0 {
			var buf bytes.Buffer
			buf.WriteRune('[')
			result = append(result, buf.Bytes()...)
		}
		s, _ := p.marshalJSON()
		result = append(result, s...)
		if index == total-1 {
			var buf bytes.Buffer
			buf.WriteRune(']')
			result = append(result, buf.Bytes()...)
		} else {
			var buf bytes.Buffer
			buf.WriteRune(',')
			result = append(result, buf.Bytes()...)
		}
	}

	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(result)
}

// MonitorStartGet starts to monior the input api
func MonitorStartGet(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")
	api, fetchErr := model.APIByID(apiID)
	if fetchErr != nil {
		log.Println(fetchErr)
		sess.AddFlash(view.Flash{Message: "接口不存在!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		startErr := service.StartMonitor(api)
		if startErr != nil {
			sess.AddFlash(view.Flash{Message: "监控启动失败!", Class: view.FlashSuccess})
			sess.Save(r, w)
			return
		}
		sess.AddFlash(view.Flash{Message: "监控启动成功!", Class: view.FlashSuccess})
		sess.Save(r, w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

// MonitorPauseGet starts to monior the input api
func MonitorPauseGet(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")
	api, fetchErr := model.APIByID(apiID)
	if fetchErr != nil {
		log.Println(fetchErr)
		sess.AddFlash(view.Flash{Message: "接口不存在!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		startErr := service.PauseMonitor(api)
		if startErr != nil {
			sess.AddFlash(view.Flash{Message: "监控暂停失败!", Class: view.FlashSuccess})
			sess.Save(r, w)
			return
		}
		sess.AddFlash(view.Flash{Message: "监控暂停!", Class: view.FlashSuccess})
		sess.Save(r, w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
