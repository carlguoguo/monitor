package model

import (
	"app/shared/database"
	"time"
)

// *****************************************************************************
// API
// *****************************************************************************

// API table contains the information for each api
type API struct {
	ID           uint32    `db:"id"`
	URL          string    `db:"url"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatedBy    uint32    `db:"user_id"`
	IntervalTime uint8     `db:"interval_time"`
}

// Request table contains the information for each request status for specific api
type Request struct {
	ID     uint32    `db:"id"`
	APIID  uint32    `db:"api_id"`
	Status string    `db:"status"`
	Time   time.Time `db:"time"`
}

// APIs return all api
func APIs() ([]API, error) {
	var result []API
	err := database.SQL.Select(&result, "SELECT id, url, interval_time, user_id, created_at, updated_at FROM api")
	return result, standardizeError(err)
}

// APIByID gets api by ID
func APIByID(apiID string) (API, error) {
	result := API{}
	err := database.SQL.Get(&result, "SELECT id, url, interval_time, user_id, created_at, updated_at FROM api WHERE id = ? LIMIT 1", apiID)
	return result, standardizeError(err)
}

// APICreate creates an api
func APICreate(url string, intervalTime int, userID string) error {
	_, err := database.SQL.Exec("INSERT INTO api (url, interval_time, user_id) VALUES (?,?,?)", url, intervalTime, userID)
	return standardizeError(err)

}

// APIUpdate updates an api
func APIUpdate(url string, intervalTime int, userID string, apiID string) error {
	_, err := database.SQL.Exec("UPDATE api SET url=?, interval_time=?, user_id=? WHERE id=? LIMIT 1", url, intervalTime, userID, apiID)
	return standardizeError(err)
}

// APIDelete deletes a note
func APIDelete(apiID string) error {
	_, err := database.SQL.Exec("DELETE FROM api WHERE id = ?", apiID)
	return standardizeError(err)
}
