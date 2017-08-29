package model

import (
	"app/shared/database"
)

// *****************************************************************************
// APIStatus
// *****************************************************************************

// APIStatus table contains the information for each api
type APIStatus struct {
	ID                  uint32  `db:"id"`
	APIID               uint32  `db:"api_id"`
	Status              int     `db:"status"`
	UpPercentage        float64 `db:"up_percentage"`
	UpSince             int     `db:"up_since"`
	AverageResponseTime int     `db:"average_response_time"`
	Count               int     `db:"count"`
	OKCount             int     `db:"ok_count"`
}

// APIStatusCreate creates an api status
func APIStatusCreate(apiID string) error {
	_, err := database.SQL.Exec("INSERT INTO api_status (api_id) VALUES (?)", apiID)
	return standardizeError(err)
}

// APIStatusUpdate updates an api status
func APIStatusUpdate(apiID string, status int, count int, okCount int, upPercentage float64, averageResponseTime int) error {
	_, err := database.SQL.Exec("UPDATE api_status SET count=?, ok_count=?, status=?, up_percentage=?, average_response_time=? WHERE api_id=? LIMIT 1", count, okCount, status, upPercentage, averageResponseTime, apiID)
	return standardizeError(err)
}

// APIStatusByID find the api status by api id
func APIStatusByID(apiID string) (APIStatus, error) {
	result := APIStatus{}
	err := database.SQL.Get(&result, "SELECT * FROM api_status WHERE api_id = ? LIMIT 1", apiID)
	return result, standardizeError(err)
}

// APIStatusAll find the api status by api id
func APIStatusAll() ([]APIStatus, error) {
	var result []APIStatus
	err := database.SQL.Select(&result, "SELECT * FROM api_status")
	return result, standardizeError(err)
}

// APIStatusUpdateAndReturn updates an api status and return it
func APIStatusUpdateAndReturn(apiID string, status int, count int, okCount int, upPercentage float64, averageResponseTime int) (APIStatus, error) {
	var err error
	updateErr := APIStatusUpdate(apiID, status, count, okCount, upPercentage, averageResponseTime)
	if updateErr != nil {
		err = updateErr
		return APIStatus{}, standardizeError(err)
	}
	apiStatus, fetchErr := APIStatusByID(apiID)
	if fetchErr != nil {
		err = fetchErr
		return APIStatus{}, standardizeError(err)
	}
	return apiStatus, standardizeError(err)
}
