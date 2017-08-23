package model

import (
	"app/shared/database"
	"time"
)

// Request table contains the information for each request status for specific api
type Request struct {
	ID          uint32    `db:"id"`
	APIID       uint32    `db:"api_id"`
	Status      int       `db:"status"`
	Cost        int       `db:"cost"`
	ContentSize int       `db:"content_size"`
	RequestTime time.Time `db:"request_time"`
}

// RequestCreate creates an api
func RequestCreate(apiID string, status int, cost int, contentSize int) error {
	_, err := database.SQL.Exec("INSERT INTO request (api_id, status, cost, content_size) VALUES (?,?,?,?)",
		apiID, status, cost, contentSize)
	return standardizeError(err)
}

// RequestByAPIID get all requests by api id
func RequestByAPIID(apiID string) ([]Request, error) {
	var result []Request
	err := database.SQL.Select(&result, "SELECT * FROM request WHERE api_id = ?", apiID)
	return result, standardizeError(err)
}
