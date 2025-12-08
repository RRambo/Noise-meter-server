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

// PostHandler handles real-time sound data from Arduino
func PostHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
    var data models.Data

    // Decode JSON payload
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
        return
    }

    // Fill timestamp if missing
    if data.MeasureTime == "" {
        data.MeasureTime = time.Now().Format(time.RFC3339)
    }

    // Fill default device ID if missing
    if data.DeviceID == "" {
        data.DeviceID = "arduino_001"
    }

    // Fill default threshold if missing
    if data.Threshold == 0 {
        data.Threshold = 70
    }

    // === Vertical logging with all required fields ===
    logMsg := "\n------ Measurement ------\n" +
        fmt.Sprintf("ID              : %d\n", data.ID) +
        fmt.Sprintf("Device ID       : %s\n", data.DeviceID) +
        fmt.Sprintf("Room Name       : %s\n", data.RoomName) +
        fmt.Sprintf("Measure Time    : %s\n", data.MeasureTime) +
        fmt.Sprintf("Is Periodic     : %t\n", data.IsPeriodic) +
        fmt.Sprintf("Sound Level     : %.2f dB\n", data.SoundLevel) +
        fmt.Sprintf("Threshold       : %.2f dB\n", data.Threshold) +
        fmt.Sprintf("5-min Average   : %.2f dB\n", data.FiveMinuteAverageSoundLevel) +
        fmt.Sprintf("LED Status      : %s\n", data.LedStatus) +
        "-------------------------"

    logger.Println(logMsg)

    // Context with timeout
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    // Save depending on periodic flag
    if !data.IsPeriodic {
        if err := ds.CreateLatest(&data, ctx); err != nil {
            switch err.(type) {
            case service.DataError:
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(`{"error": "` + err.Error() + `"}`))
                return
            default:
                logger.Println("Error updating latest measurement data:", err, data)
                http.Error(w, "Internal server error.", http.StatusInternalServerError)
                return
            }
        }
    } else {
        if err := ds.Create(&data, ctx); err != nil {
            switch err.(type) {
            case service.DataError:
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(`{"error": "` + err.Error() + `"}`))
                return
            default:
                logger.Println("Error creating data:", err, data)
                http.Error(w, "Internal server error.", http.StatusInternalServerError)
                return
            }
        }
    }

    // Return the created record as JSON
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        logger.Println("Error encoding data:", err, data)
        http.Error(w, "Internal server error.", http.StatusInternalServerError)
        return
    }
}