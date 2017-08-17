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
		apis, err := model.APIs()
		if err != nil {
			log.Println(err)
			apis = []model.API{}
		}

		// Display the view
		v := view.New(r)
		v.Name = "index/auth"
		v.Vars["username"] = session.Values["username"]
		v.Vars["apis"] = apis
		v.Render(w)
	} else {
		// Display the view
		v := view.New(r)
		v.Name = "index/anon"
		v.Render(w)
		return
	}
}
