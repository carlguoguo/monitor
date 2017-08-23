package route

import (
	"net/http"

	"app/controller"
	"app/route/middleware/acl"
	hr "app/route/middleware/httprouterwrapper"
	"app/route/middleware/logrequest"
	"app/route/middleware/pprofhandler"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// LoadHTTP returns the HTTP routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())
}

// LoadHTTPS returns the HTTPS routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************
func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page
	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.IndexGET)))

	// Login
	r.GET("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
	r.GET("/logout", hr.Handler(alice.
		New().
		ThenFunc(controller.LogoutGET)))

	// Register
	r.GET("/register", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterGET)))
	r.POST("/register", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterPOST)))

	// API
	r.POST("/api/create", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.APICreatePOST)))
	r.POST("/api/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.APIUpdatePost)))
	r.GET("/api/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.APIDeleteGet)))
	r.GET("/api/detail/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.APIDetailGet)))

	// Montor
	r.GET("/monitor/start/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MonitorStartGet)))
	r.GET("/monitor/pause/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MonitorPauseGet)))

	// Enable Pprof
	r.GET("/debug/pprof/*pprof", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(pprofhandler.Handler)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************
func middleware(h http.Handler) http.Handler {
	// Log every Request
	h = logrequest.Handler(h)
	return h
}
