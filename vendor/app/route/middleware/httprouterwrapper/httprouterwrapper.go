// Package httprouterwrapper allows the use of http.HandlerFunc compatible funcs with julienschmidt/httprouter
package httprouterwrapper

import (
    "net/http"

    "github.com/gorilla/context"
    "github.com/julienschmidt/httprouter"
)

// Source: http://nicolasmerouze.com/guide-routers-golang/

func Handler(h http.Handler) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        context.Set(r, "params", p)
        h.ServeHTTP(w, r)
    }
}
