package model

import (
	"time"

	"app/shared/database"
	"app/shared/interval"
)

// *****************************************************************************
// API
// *****************************************************************************

// API table contains the information for each api
type API struct {
	ID             uint32    `db:"id"`
	Alias          string    `db:"alias"`
	URL            string    `db:"url"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	AlertReceivers string    `db:"alert_receivers"`
	Timeout        int       `db:"timeout"`
	FailMax        uint8     `db:"fail_max"`
	CreatedBy      uint32    `db:"user_id"`
	IntervalTime   uint8     `db:"interval_time"`
	Start          uint8     `db:"start"`
	Job            interval.Job
}

// APIs return all api
func APIs() ([]API, error) {
	var result []API
	err := database.SQL.Select(&result, "SELECT id, start, url, interval_time, user_id, alias, alert_receivers, timeout, fail_max, created_at, updated_at FROM api")
	return result, standardizeError(err)
}

// StartedAPIs return all started api
func StartedAPIs() ([]API, error) {
	var result []API
	err := database.SQL.Select(&result, "SELECT id, start, url, interval_time, user_id, alias, alert_receivers, timeout, fail_max, created_at, updated_at FROM api where start = 1")
	return result, standardizeError(err)
}

// APIByID gets api by ID
func APIByID(apiID string) (API, error) {
	result := API{}
	err := database.SQL.Get(&result, "SELECT id, start, url, interval_time, user_id, alias, alert_receivers, timeout, fail_max, created_at, updated_at FROM api WHERE id = ? LIMIT 1", apiID)
	return result, standardizeError(err)
}

// APIByURL gets api by ID
func APIByURL(url string) (API, error) {
	result := API{}
	err := database.SQL.Get(&result, "SELECT id, start, url, interval_time, user_id, alias, alert_receivers, timeout, fail_max, created_at, updated_at FROM api WHERE url = ? LIMIT 1", url)
	return result, standardizeError(err)
}

// APICreate creates an api
func APICreate(url string, intervalTime int, userID string, alias string, alertReceivers string, timeout int, failMax int) error {
	_, err := database.SQL.Exec("INSERT INTO api (url, interval_time, user_id, alias, alert_receivers, timeout, fail_max) VALUES (?,?,?,?,?,?,?)", url, intervalTime, userID, alias, alertReceivers, timeout, failMax)
	return standardizeError(err)
}

// APICreateAndGet creates an api and return API struct
func APICreateAndGet(url string, intervalTime int, userID string, alias string, alertReceivers string, timeout int, failMax int) (API, error) {
	var err error
	creationErr := APICreate(url, intervalTime, userID, alias, alertReceivers, timeout, failMax)
	if creationErr != nil {
		err = creationErr
		return API{}, standardizeError(err)
	}
	api, fetchErr := APIByURL(url)
	if fetchErr != nil {
		err = fetchErr
	}
	return api, standardizeError(err)
}

// APIUpdate updates an api
func APIUpdate(url string, intervalTime int, userID string, alias string, alertReceivers string, timeout int, failMax int, apiID string) error {
	_, err := database.SQL.Exec("UPDATE api SET url=?, start=0, interval_time=?, user_id=?, alias=?, alert_receivers=?, timeout=?, fail_max=? WHERE id=? LIMIT 1", url, intervalTime, userID, alias, alertReceivers, timeout, failMax, apiID)
	return standardizeError(err)
}

// APIUpdateStart updates an api
func APIUpdateStart(start int, apiID string) error {
	_, err := database.SQL.Exec("UPDATE api SET start=? WHERE id=? LIMIT 1", start, apiID)
	return standardizeError(err)
}

// APIUpdateAndReturn updates an api and return it
func APIUpdateAndReturn(url string, intervalTime int, userID string, alias string, alertReceivers string, timeout int, failMax int, apiID string) (API, error) {
	var err error
	updateErr := APIUpdate(url, intervalTime, userID, alias, alertReceivers, timeout, failMax, apiID)
	if updateErr != nil {
		err = updateErr
	}
	api, fetchErr := APIByURL(url)
	if fetchErr != nil {
		err = fetchErr
	}
	return api, standardizeError(err)
}

// APIDelete deletes an api
func APIDelete(apiID string) error {
	_, err := database.SQL.Exec("DELETE FROM api WHERE id = ?", apiID)
	return standardizeError(err)
}

// APIDeleteAndReturn deletes an api and return
func APIDeleteAndReturn(apiID string) (API, error) {
	var err error
	api, fetchErr := APIByID(apiID)
	if fetchErr != nil {
		err = fetchErr
	}
	deleteErr := APIDelete(apiID)
	if deleteErr != nil {
		err = deleteErr
	}
	return api, standardizeError(err)
}
