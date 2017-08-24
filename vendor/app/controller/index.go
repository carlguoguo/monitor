package controller

import (
	"log"
	"net/http"

	"app/model"
	"app/shared/session"
	"app/shared/view"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := session.Instance(r)
	if session.Values["id"] != nil {
		apisStatus := make(map[uint32]model.APIStatus)
		apis, err := model.APIs()
		if err != nil {
			log.Println(err)
			apis = []model.API{}
		} else {
			apiStatusAll, err := model.APIStatusAll()
			if err != nil {
				log.Println(err)
			} else {
				for _, apiStatus := range apiStatusAll {
					apisStatus[apiStatus.APIID] = apiStatus
				}
			}
		}
		// Display the view
		v := view.New(r)
		v.Name = "index/auth"
		v.Vars["username"] = session.Values["username"]
		v.Vars["apis"] = apis
		v.Vars["apisStatus"] = apisStatus
		v.Render(w)
	} else {
		LoginGET(w, r)
	}
}
