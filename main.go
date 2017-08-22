package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"app/model"
	"app/route"
	"app/service"
	"app/shared/database"
	"app/shared/email"
	"app/shared/jsonconfig"
	"app/shared/server"
	"app/shared/session"
	"app/shared/view"
	"app/shared/view/plugin"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	email.Configure(config.Email)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
	)

	// Start Interval Job
	apis, err := model.APIs()
	if err != nil {
		log.Fatal(err)
	}
	for _, api := range apis {
		service.NewMonitor(api)
		if api.Start == 1 {
			service.StartMonitor(api)
		}
	}

	// Start the listener
	server.Run(route.LoadHTTP(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database database.MySQLInfo `json:"Database"`
	Email    email.SMTPInfo     `json:"Email"`
	Server   server.Server      `json:"Server"`
	Session  session.Session    `json:"Session"`
	Template view.Template      `json:"Template"`
	View     view.View          `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
