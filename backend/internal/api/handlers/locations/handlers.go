package locations

import (
	"encoding/json"
	"goapi/internal/api/repository/models"
	"log"
	"net/http"
	"strconv"
)

type LocationService interface {
	CreateLocation(location *models.Location) error
	GetAllLocations() ([]*models.Location, error)
	SetChosenLocation(id int) error
}

func GetLocationsHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc LocationService) {
	locations, err := svc.GetAllLocations()
	if err != nil {
		logger.Println("Error getting locations:", err)
		http.Error(w, `{"error": "Failed to get locations"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"locations": locations,
	})
}

func CreateLocationHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc LocationService) {
	var location models.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := svc.CreateLocation(&location); err != nil {
		logger.Println("Error creating location:", err)
		http.Error(w, `{"error": "Failed to create location"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}

func SetChosenLocationHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc LocationService) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid location ID"}`, http.StatusBadRequest)
		return
	}

	if err := svc.SetChosenLocation(id); err != nil {
		logger.Println("Error setting chosen location:", err)
		http.Error(w, `{"error": "Failed to set chosen location"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Chosen location updated"}`))
}
