package data

import (
    "context"
    "encoding/json"
    "fmt"
    "goapi/internal/api/repository/models"
    service "goapi/internal/api/service/data"
    "log"
    "net/http"
    "time"
)

// PostHourlyHandler handles hourly and real-time average sound data from Arduino
func PostHourlyHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
    var data models.Data

    // Decode JSON payload
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        logger.Println("❌ Failed to decode JSON payload:", err)
        http.Error(w, `{"error": "Invalid request data. Please check your input."}`, http.StatusBadRequest)
        return
    }

    // Fill missing timestamp (SQLite-friendly format)
    if data.MeasureTime == "" {
        data.MeasureTime = time.Now().UTC().Format("2006-01-02 15:04:05")
    }

    // Fill missing device ID
    if data.DeviceID == "" {
        data.DeviceID = "arduino_001"
    }

    // Pretty measurement output (no logger prefix)
    logMsg := fmt.Sprintf(`
------ Measurement ------
ID              : %d
Device ID       : %s
Room Name       : %s
Measure Time    : %s
Is Periodic     : %t
Sound Level     : %.2f dB
Threshold       : %.2f dB
5-min Average   : %.2f dB
LED Status      : %s
-------------------------`,
        data.ID,
        data.DeviceID,
        data.RoomName,
        data.MeasureTime,
        data.IsPeriodic,
        data.SoundLevel,
        data.Threshold,
        data.FiveMinuteAverageSoundLevel,
        data.LedStatus,
    )

    fmt.Println(logMsg)

    // Context with timeout
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    // Save to repository
    if err := ds.Create(&data, ctx); err != nil {
        switch err.(type) {
        case service.DataError:
            logger.Println("⚠️ Data validation error:", err)
            http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
            return
        default:
            logger.Println("❌ Internal error creating hourly data:", err)
            http.Error(w, "Internal server error.", http.StatusInternalServerError)
            return
        }
    }

    // Return the created record as JSON
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        logger.Println("❌ Error encoding response JSON:", err)
        http.Error(w, "Internal server error.", http.StatusInternalServerError)
        return
    }
}
