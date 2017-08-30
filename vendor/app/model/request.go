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
func RequestByAPIID(apiID string, limit int) ([]Request, error) {
	var result []Request
	var err error
	query := "SELECT * FROM request WHERE api_id = ? order by request_time desc"
	if limit > 0 {
		err = database.SQL.Select(&result, query+" limit ?", apiID, limit)
	}
	return result, standardizeError(err)
}
