package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

var (
	// SQL Connection
	SQL *sqlx.DB
	// Database Info
	databases MySQLInfo
)

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// DSN returns the Data Source Name
func DSN(ci MySQLInfo) string {
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

// Connect to the database
func Connect(d MySQLInfo) {
	var err error

	// Store the config
	databases = d
	// Connect to MySQL
	if SQL, err = sqlx.Connect("mysql", DSN(d)); err != nil {
		log.Println("SQL Driver Error", err)
	}

	// Check if is alive
	if err = SQL.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}

// ReadConfig returns the database information
func ReadConfig() MySQLInfo {
	return databases
}
