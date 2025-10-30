package data

import (
	"context"
	"goapi/internal/api/repository/models"
)

// * Implementation of DataService for SQLite database *
type DataServiceSQLite struct {
	repo models.DataRepository
}

func NewDataServiceSQLite(repo models.DataRepository) *DataServiceSQLite {
	return &DataServiceSQLite{
		repo: repo,
	}
}

func (ds *DataServiceSQLite) Create(data *models.Data, ctx context.Context) error {

	if err := ds.ValidateData(data); err != nil {
		return DataError{Message: "Invalid data: " + err.Error()}
		// There was bug here, now it returns error massage.
		// Original line : return DataError{Message: "InvalMockDataServiceSuccessfulid data."}
	}
	return ds.repo.Create(data, ctx)
}

func (ds *DataServiceSQLite) ReadOne(id int, ctx context.Context) (*models.Data, error) {

	data, err := ds.repo.ReadOne(id, ctx)
	if err != nil {
		return nil, err
	}

	_ = data

	// Tehdään datalle jotain, päätellään datasta jotain!!!
	// Tämä ohjaa toimintaa älykkäästi, esim. jos data on tietynlaista, niin tehdään jotain
	// (What do these mean??)

	// We do something to the data, we deduce something from the data!!!
	// This guides the operation intelligently, for example, if the data is of a certain type, then we do something

	return data, nil
}

func (ds *DataServiceSQLite) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return ds.repo.ReadMany(page, rowsPerPage, ctx)
}

func (ds *DataServiceSQLite) Update(data *models.Data, ctx context.Context) (int64, error) {

	if err := ds.ValidateData(data); err != nil {
		return 0, DataError{Message: "Invalid data: " + err.Error()}
	}
	return ds.repo.Update(data, ctx)
}

func (ds *DataServiceSQLite) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return ds.repo.Delete(data, ctx)
}

func (ds *DataServiceSQLite) ValidateData(data *models.Data) error {
	var errMsg string
	if data.DeviceID == "" || len(data.DeviceID) > 50 {
		errMsg += "DeviceID is required and must be less than 50 characters. "
	}
	if data.RoomName == "" {
		errMsg += "RoomName is required. "
	}
	// Maybe we need to edit the system around RoomName so typos don't messup the data, for robustness and better usability.
	// Maybe by making a predefined list of room names to choose from in the frontend.
	if data.SoundLevel < 0 || data.SoundLevel > 150 {
		errMsg += "SoundLevel must be between 0 and 150 dB. "
	}
	if data.Threshold < 0 || data.Threshold > 150 {
		errMsg += "Threshold must be between 0 and 150 dB. "
	}

	if errMsg != "" {
		return DataError{Message: errMsg}
	}
	// Apparently there was bug here too, it didn't return error message when there was validation error.

	return nil
}
