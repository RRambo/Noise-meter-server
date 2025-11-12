package data

import (
	"context"
	"encoding/json"
	"goapi/internal/api/repository/models"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	var data models.Data

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	// Fill timestamp if missing
	if data.MeasureTime == "" {
		data.MeasureTime = time.Now().Format(time.RFC3339)
	}
	if data.DeviceID == "" {
		data.DeviceID = "arduino_001"
	}

	logger.Println("Received POST /data from Arduino:")
	logger.Printf("%+v\n", data)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

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

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Println("Error encoding data:", err, data)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
