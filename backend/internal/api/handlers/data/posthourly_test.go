package data_test

import (
	"encoding/json"
	"goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostHourlyHandlerSuccess(t *testing.T) {
	mockDS := &service.MockDataServiceSuccessful{}

	jsonPayload := `{
		"device_id": "arduino_001",
		"room_name": "Class_A",
		"hourly_average_sound_level": 63.45,
		"description": "1-hour average sound reading"
	}`

	req, err := http.NewRequest("POST", "/data/hourly", strings.NewReader(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PostHourlyHandler(rr, req, log.Default(), mockDS)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var respData map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&respData); err != nil {
		t.Errorf("could not decode response: %v", err)
	}
	if respData["device_id"] != "arduino_001" {
		t.Errorf("unexpected device_id: %v", respData["device_id"])
	}
}