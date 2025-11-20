package data_test

import (
	"encoding/json"
	"goapi/internal/api/handlers/data"
	"goapi/internal/api/repository/models"
	service "goapi/internal/api/service/data"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostInvalidRequestBody(t *testing.T) {

	req, err := http.NewRequest("POST", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Body = io.NopCloser(strings.NewReader(`Plain text, not JSON`))
	rr := httptest.NewRecorder()
	data.PostHandler(rr, req, log.Default(), &service.MockDataServiceSuccessful{})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Invalid request data. Please check your input."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPostErrorCreatingData(t *testing.T) {

	req, err := http.NewRequest("POST", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	dataJSON, _ := json.Marshal(models.Data{
		ID:          1,
		DeviceID:    "arduino wifi R2",
		RoomName:    "eating room",
		SoundLevel:  65.5,
		Threshold:   70.0,
		MeasureTime: "2024-06-01T12:00:00Z",
		IsAlert:     false,
		Description: "post test data 1",
	})

	req.Body = io.NopCloser(strings.NewReader(string(dataJSON)))
	rr := httptest.NewRecorder()

	data.PostHandler(rr, req, log.Default(), &service.MockDataServiceError{})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Error creating data."}` // * This message is passed from the MockDataServiceError
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPostSuccessful(t *testing.T) {

	req, err := http.NewRequest("POST", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	dataJSON, _ := json.Marshal(models.Data{
		ID:          1,
		DeviceID:    "arduino wifi R2",
		RoomName:    "eating room",
		SoundLevel:  65.5,
		Threshold:   70.0,
		MeasureTime: "2024-06-01T12:00:00Z",
		IsAlert:     false,
		Description: "post test data 1",
	})

	// * Create new reader with the JSON payload
	req.Body = io.NopCloser(strings.NewReader(string(dataJSON)))

	rr := httptest.NewRecorder()

	// * Call the handler
	data.PostHandler(rr, req, log.Default(), &service.MockDataServiceSuccessful{})

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// * Check the response body
	expected := `{"id":1,"device_id":"arduino wifi R2","room_name":"eating room","sound_level":65.5,"threshold":70,"measure_time":"2024-06-01T12:00:00Z","is_alert":false,"description":"post test data 1","current_sound_level":65.5,"average_sound_level":65.5}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
