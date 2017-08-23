package controller

import (
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
		model.APIStatusCreate(newAPI.ID)
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
