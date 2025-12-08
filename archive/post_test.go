package data_test

import (
    "encoding/json"
    "fmt"
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
    mockDS := &service.MockDataServiceSuccessful{}

    req, err := http.NewRequest("POST", "/data", nil)
    if err != nil {
        t.Fatal(err)
    }

    req.Body = io.NopCloser(strings.NewReader(`Plain text, not JSON`))
    rr := httptest.NewRecorder()

    data.PostHandler(rr, req, log.Default(), mockDS)

    if rr.Code != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
    }

    expected := `{"error": "Invalid request data. Please check your input."}`
    if strings.TrimSpace(rr.Body.String()) != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestPostErrorCreatingData(t *testing.T) {
    mockDS := &service.MockDataServiceError{}

    payload := models.Data{
        ID:          1,
        DeviceID:    "arduino wifi R2",
        RoomName:    "eating room",
        SoundLevel:  65.5,
        Threshold:   70.0,
        MeasureTime: "2024-06-01T12:00:00Z",
        IsPeriodic:  true,
        FiveMinuteAverageSoundLevel: 64.8,
        LedStatus:   "OFF",
    }
    dataJSON, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", "/data", io.NopCloser(strings.NewReader(string(dataJSON))))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    data.PostHandler(rr, req, log.Default(), mockDS)

    if rr.Code != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
    }

    expected := `{"error": "Error creating data."}`
    if strings.TrimSpace(rr.Body.String()) != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestPostSuccessful(t *testing.T) {
    mockDS := &service.MockDataServiceSuccessful{}

    payload := models.Data{
        ID:          1,
        DeviceID:    "arduino wifi R2",
        RoomName:    "eating room",
        SoundLevel:  65.5,
        Threshold:   70.0,
        MeasureTime: "2024-06-01T12:00:00Z",
        IsPeriodic:  true,
        FiveMinuteAverageSoundLevel: 64.8,
        LedStatus:   "OFF",
    }
    dataJSON, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", "/data", io.NopCloser(strings.NewReader(string(dataJSON))))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    data.PostHandler(rr, req, log.Default(), mockDS)

    if rr.Code != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
    }

    expected := fmt.Sprintf(`{"id":1,"device_id":"%s","room_name":"%s","sound_level":65.5,"threshold":70,"measure_time":"%s","is_periodic":true,"five_minute_average_sound_level":64.8,"led_status":"OFF"}`,
        payload.DeviceID, payload.RoomName, payload.MeasureTime)

    if strings.TrimSpace(rr.Body.String()) != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}
