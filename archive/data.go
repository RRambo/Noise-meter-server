package models

// Though this file is named data, it contains only the sound data.
// At this point I'll leave it like this. If we rename it it will require far more changes on other files.

import (
	"context"
	"time"
)

type Data struct {
	ID          int     `json:"id,omitempty"`
	DeviceID    string  `json:"device_id"`    // Arduino device ID
	RoomName    string  `json:"room_name"`    // Name of the working room
	SoundLevel  float64 `json:"sound_level"`  // Level of sound in dB
	Threshold   float64 `json:"threshold"`    // Threshold level in dB
	MeasureTime string  `json:"measure_time"` // Time of measurement
    LedStatus   string  `json:"led_status,omitempty"`   // Actual LED state from Arduino	// Need to figure out how to make this be calculated and filled by the soundlevel&threshold system
	Description             string  `json:"description"`                          // Additional information
	IsPeriodic              bool    `json:"is_periodic,omitempty"`                // Is the data constantly/periodically measured
	HourlyAverageSoundLevel float64 `json:"hourly_average_sound_level,omitempty"` // hourly average
	FiveMinuteAverageSoundLevel float64 `json:"five_minute_average_sound_level,omitempty"` // 5-minute average
}

type DataRepository interface {
    Create(data *Data, ctx context.Context) error
    CreateLatest(data *Data, ctx context.Context) error
    ReadOne(id int, ctx context.Context) (*Data, error)
    ReadLatest(id string, ctx context.Context) (*Data, error)
    ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*Data, error)
    Update(data *Data, ctx context.Context) (int64, error)
    Delete(data *Data, ctx context.Context) (int64, error)
    GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*Data, error)
    GetByRoom(roomName string, ctx context.Context) ([]*Data, error)
    ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error)
}
