package SQLite

import (
	"context"
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/models"
	"time"
)

type DataRepository struct {
	sqlDB *sql.DB
	createStmt,
	readStmt,
	readManyStmt,
	updateStmt,
	deleteStmt *sql.Stmt
	ctx context.Context
}

func NewDataRepository(sqlDB DAL.SQLDatabase, ctx context.Context) (models.DataRepository, error) {

	repo := &DataRepository{
		sqlDB: sqlDB.Connection(),
		ctx:   ctx,
	}

	// Create the data table if it doesn't exist
	if _, err := repo.sqlDB.Exec(`CREATE TABLE  IF NOT EXISTS data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id TEXT NOT NULL DEFAULT 'arduino_001',
		room_name TEXT NOT NULL DEFAULT 'unassigned',
		sound_level REAL NOT NULL DEFAULT 0.0,
		threshold REAL NOT NULL DEFAULT 70.0, 
		measure_time TEXT NOT NULL,
		is_alert INTEGER NOT NULL DEFAULT 0,
		description TEXT DEFAULT ''
	);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// * Create needed Prepared SQL statements, this is more efficient than running each query individually
	createStmt, err := repo.sqlDB.Prepare(`INSERT INTO data (
	device_id, room_name, sound_level, threshold, measure_time, is_alert, description)
	VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		repo.sqlDB.Close() // Close the database connection if statement preparation fails
		return nil, err
	}
	repo.createStmt = createStmt

	// Read single record
	readStmt, err := repo.sqlDB.Prepare(`SELECT 
	id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description
	FROM data WHERE id = ?`) // I think there was bug here with mistaked rune ("" instead of ``)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readStmt = readStmt

	// Read multiple records with pagination
	readManyStmt, err := repo.sqlDB.Prepare(`SELECT
	id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description
	FROM data LIMIT ? OFFSET ?`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readManyStmt = readManyStmt

	// Update record
	updateStmt, err := repo.sqlDB.Prepare(`UPDATE data SET 
	device_id = ?, room_name = ?, sound_level = ?, threshold = ?, 
	measure_time = ?, is_alert = ?, description = ?  
	WHERE id = ?`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.updateStmt = updateStmt

	// Delete record
	deleteStmt, err := repo.sqlDB.Prepare(`DELETE FROM data WHERE id = ?`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.deleteStmt = deleteStmt

	go Close(ctx, repo)

	return repo, nil
}

func Close(ctx context.Context, r *DataRepository) {

	<-ctx.Done()
	r.createStmt.Close()
	r.readStmt.Close()
	r.updateStmt.Close()
	r.deleteStmt.Close()
	r.readManyStmt.Close()
	r.sqlDB.Close()
}

func (r *DataRepository) Create(data *models.Data, ctx context.Context) error {

	// Set default values if not provided
	if data.DeviceID == "" {
		data.DeviceID = "arduino_001"
	}
	if data.Threshold == 0 {
		data.Threshold = 70.0
	}
	if data.SoundLevel >= data.Threshold {
		data.IsAlert = true
	}
	if data.MeasureTime == "" {
		data.MeasureTime = time.Now().Format(time.RFC3339)
		//Fill with current time if not provided
	}

	// Execute INSERT with correct field order
	res, err := r.createStmt.ExecContext(ctx,
		data.DeviceID,
		data.RoomName,
		data.SoundLevel,
		data.Threshold,
		data.MeasureTime,
		data.IsAlert,
		data.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	data.ID = int(id)
	return nil
}

func (r *DataRepository) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	row := r.readStmt.QueryRowContext(ctx, id)
	var data models.Data

	// Scan with correct field order matching new schema
	err := row.Scan(
		&data.ID,
		&data.DeviceID,
		&data.RoomName,
		&data.SoundLevel,
		&data.Threshold,
		&data.MeasureTime,
		&data.IsAlert,
		&data.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *DataRepository) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	if page < 1 {
		return r.ReadAll()
	}

	offset := rowsPerPage * (page - 1)
	rows, err := r.readManyStmt.QueryContext(ctx, rowsPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (r *DataRepository) ReadAll() ([]*models.Data, error) {
	rows, err := r.sqlDB.QueryContext(context.Background(),
		`SELECT id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description 
		FROM data`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (r *DataRepository) Update(data *models.Data, ctx context.Context) (int64, error) {
	// Update measure_time if modifying
	if data.MeasureTime == "" {
		data.MeasureTime = time.Now().Format(time.RFC3339)
	}

	res, err := r.updateStmt.ExecContext(ctx,
		data.DeviceID,
		data.RoomName,
		data.SoundLevel,
		data.Threshold,
		data.MeasureTime,
		data.IsAlert,
		data.Description,
		data.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *DataRepository) Delete(data *models.Data, ctx context.Context) (int64, error) {
	res, err := r.deleteStmt.ExecContext(ctx, data.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *DataRepository) GetByRoom(roomName string, ctx context.Context) ([]*models.Data, error) {
	// Example implementation
	rows, err := r.sqlDB.QueryContext(ctx, `SELECT id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description FROM data WHERE room_name = ?`, roomName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (r *DataRepository) GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	// Basic stub: returns all data for the room (replace with real daily summary logic as needed)

	// Calculate start and end of the day
	year, month, day := date.Date()
	location := date.Location()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, location)
	endOfDay := startOfDay.Add(24 * time.Hour)

	startOfDayStr := startOfDay.UTC().Format(time.RFC3339)
	endOfDayStr := endOfDay.UTC().Format(time.RFC3339)

	// Query to get data for the specified room and date range
	query := `
	SELECT id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description 
	FROM data
	WHERE room_name = ?
		AND measure_time >= ?
		AND measure_time < ?
	ORDER BY measure_time ASC
	`
	rows, err := r.sqlDB.QueryContext(ctx, query, roomName, startOfDayStr, endOfDayStr)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}
