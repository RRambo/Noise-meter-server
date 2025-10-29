package data_test

import (
	"goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPutInvalidRequestBody(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Body = io.NopCloser(strings.NewReader(`Plain text, not JSON`))
	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceSuccessful{})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Invalid request data. Please check your input."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPutHandlerError(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", strings.NewReader(`{"id": 1, "device_id": "arduino_001", "room_name": "PlayRoom_A", "sound_level": 75.5, "threshold": 70.0, "measure_time": "2024-01-01T12:00:00Z", "is_alert": true, "description": "Test update"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceError{})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Error updating data."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPutDataNotFound(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", strings.NewReader(`{"id": 999, "device_id": "arduino_001", "room_name": "PlayRoom_A", "sound_level": 75.5, "threshold": 70.0, "measure_time": "2024-01-01T12:00:00Z", "is_alert": false, "description": "Test update"}`))
	// id:999 -> unexisting id for not found test
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceNotFound{})

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"error": "Resource not found."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPutHandlerSuccess(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", strings.NewReader(`{"id": 1, "device_id": "arduino_001", "room_name": "PlayRoom_B", "sound_level": 82.3, "threshold": 70.0, "measure_time": "2024-10-27T14:30:00Z", "is_alert": true, "description": "Success test"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceSuccessful{})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":1,"device_id":"arduino_001","room_name":"PlayRoom_B","sound_level":82.3,"threshold":70,"measure_time":"2024-10-27T14:30:00Z","is_alert":true,"description":"Success test"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
