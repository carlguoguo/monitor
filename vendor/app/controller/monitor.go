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

// APICreateGet displays the api creation page
func APICreateGet(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "api/create"
	v.Render(w)
}

// APICreatePOST handles the api creation form submission
func APICreatePOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	if validate, missingField := view.Validate(r, []string{"url", "interval_time"}); !validate {
		sess.AddFlash(view.Flash{Message: "缺少字段 " + missingField, Class: view.FlashError})
		sess.Save(r, w)
		APICreateGet(w, r)
		return
	}

	url := r.FormValue("url")
	intervalTime, err := strconv.Atoi(r.FormValue("interval_time"))
	if err != nil {
		sess.AddFlash(view.Flash{Message: "字段interval_time需要一个整型的数字", Class: view.FlashError})
		sess.Save(r, w)
		APICreateGet(w, r)
		return
	}

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get database result
	newAPI, creationErr := model.APICreateAndGet(url, intervalTime, userID)

	if creationErr != nil {
		log.Println(creationErr)
		sess.AddFlash(view.Flash{Message: "服务器错误，请重试!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{Message: "新接口已经添加!", Class: view.FlashSuccess})
		sess.Save(r, w)
		service.StartMonitor(newAPI)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Display the same page
	APICreateGet(w, r)
}

// APIUpdateGet displays the api update page
func APIUpdateGet(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")

	api, err := model.APIByID(apiID)
	if err != nil { // If the api doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{Message: "服务器错误，请重试！", Class: view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	v := view.New(r)
	v.Name = "api/update"
	v.Vars["url"] = api.URL
	v.Vars["interval_time"] = api.IntervalTime
	v.Render(w)
}

// APIUpdatePost handles the api update form submission
func APIUpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"url", "interval_time"}); !validate {
		sess.AddFlash(view.Flash{Message: "缺少字段: " + missingField, Class: view.FlashError})
		sess.Save(r, w)
		APIUpdateGet(w, r)
		return
	}

	// Get form values
	url := r.FormValue("url")
	intervalTime, err := strconv.Atoi(r.FormValue("interval_time"))
	if err != nil {
		sess.AddFlash(view.Flash{Message: "字段interval_time需要一个整型的数字", Class: view.FlashError})
		sess.Save(r, w)
		APICreateGet(w, r)
		return
	}

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	apiID := params.ByName("id")

	// Get database result
	newAPI, updateErr := model.APIUpdateAndReturn(url, intervalTime, userID, apiID)
	// Will only error if there is a problem with the query
	if updateErr != nil {
		log.Println(updateErr)
		sess.AddFlash(view.Flash{Message: "服务器错误，请重试!", Class: view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{Message: "接口修改成功!", Class: view.FlashSuccess})
		sess.Save(r, w)
		service.RestartMonitor(newAPI)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Display the same page
	APIUpdateGet(w, r)
}

// APIDeleteGet handles the api deletion
func APIDeleteGet(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	// Get database result
	deletedAPI, err := model.APIDeleteAndReturn(noteID)
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
