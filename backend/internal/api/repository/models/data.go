package models

// Though this file is named data, it contains only the sound data.
// At this point I'll leave it like this. If we rename it it will require far more changes on other files.

import (
	"context"
	"time"
)

type Data struct {
	ID          int     `json:"id"`
	DeviceID    string  `json:"device_id"`    // Arduino device ID
	RoomName    string  `json:"room_name"`    // Name of the working room
	SoundLevel  float64 `json:"sound_level"`  // Level of sound in dB
	Threshold   float64 `json:"threshold"`    // Threshold level in dB
	MeasureTime string  `json:"measure_time"` // Time of measurement
	IsAlert     bool    `json:"is_alert"`     // Whether the sound level exceeds the threshold
	// Need to figure out how to make this be calculated and filled by the soundlevel&threshold system
	// Done and tested but I feel like it's too rough. If needed, check internal\api\repository\DAL\SQLite\data.go line 115 for the logic fot this.
	Description string `json:"description"` // Additional information
	CurrentSoundLevel        float64 `json:"current_sound_level,omitempty"`      // Real-time
    AverageSoundLevel        float64 `json:"average_sound_level,omitempty"`      // 1-hour rolling average
    HourlyAverageSoundLevel  float64 `json:"hourly_average_sound_level,omitempty"` // hourly average
}

type DataRepository interface {
	Create(Data *Data, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*Data, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*Data, error)
	Update(data *Data, ctx context.Context) (int64, error)
	Delete(data *Data, ctx context.Context) (int64, error)
	GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*Data, error) // To retreive daily summary statistics
	GetByRoom(roomName string, ctx context.Context) ([]*Data, error)                       // (For now only for future implementation.)To retrieve data for a specific room
}
