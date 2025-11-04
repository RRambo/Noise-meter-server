# Kindergarten Sound Meter

This system supports an IoT device (Arduino with sound sensor) that measures noise levels in kindergarten rooms. When sound exceeds a threshold, data is sent to the server for logging and analysis.

This repository contains:
- A RESTful API backend for monitoring sound levels in kindergarten rooms, built with Go and SQLite.
- A frontend user interface for managing locations, configuring the noise meter, and viewing results and summaries.

### Features

- RESTful API with full CRUD operations for sound data
- SQLite database for data persistence
- Basic authentication
- Sound level validation (0-150 dB range)
- Automatic timestamping
- Alert flagging when threshold exceeded
- **Location management:** Add, select, and manage room locations via `/locations` endpoint and frontend dropdown
- **Frontend:** Simple web UI for location selection and management

## Technology Stack

- **Language:** Go 1.x
- **Database:** SQLite
- **Frontend:** HTML, CSS, JavaScript
- **Architecture:** Handler-Service-Repository pattern
- **Testing:** Go testing framework
- **Authentication:** HTTP Basic Auth

## API Endpoints

### Authentication

All endpoints require Basic Authentication:
- **Username:** `kids_noisemeter_admin`
- **Password:** `passwordkids`

### Endpoints

| Method | Endpoint         | Description                        |
|--------|------------------|------------------------------------|
| GET    | `/api/data`          | Get all sound measurements         |
| GET    | `/api/data/{id}`     | Get specific measurement           |
| POST   | `/api/data`          | Create new measurement             |
| PUT    | `/api/data`          | Update existing measurement        |
| DELETE | `/api/data/{id}`     | Delete measurement                 |
| GET    | `/api/locations`     | Get all locations                  |
| POST   | `/api/locations`     | Add a new location                 |
| PUT    | `/api/locations/{id}`| Set a location as currently chosen |

### Data Model

```json
{
  "id": 1,
  "device_id": "arduino_001",
  "room_name": "PlayRoom_A",
  "sound_level": 78.5,
  "threshold": 70.0,
  "measure_time": "2024-10-27T09:30:00Z",
  "is_alert": true,
  "description": "Morning playtime noise peak"
}
```

### Example Request

**Create Sound Measurement:**

```bash
curl -X POST http://localhost:8080/api/data \
  -u kids_noisemeter_admin:passwordkids \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "arduino_001",
    "room_name": "PlayRoom_A",
    "sound_level": 78.5,
    "threshold": 70.0,
    "measure_time": "2024-10-27T09:30:00Z",
    "is_alert": true,
    "description": "Morning playtime noise peak"
  }'
```

**Add a Location:**

```bash
curl -X POST http://localhost:8080/api/locations \
  -u kids_noisemeter_admin:passwordkids \
  -H "Content-Type: application/json" \
  -d '{ "name": "PlayRoom_A", "chosen": true }'
```

**Set Chosen Location:**

```bash
curl -X PUT http://localhost:8080/api/locations/1 \
  -u kids_noisemeter_admin:passwordkids \
  -H "Content-Type: application/json"
```

## Frontend

A simple web UI is provided in the `frontend/` folder. It allows users to:
- Add new locations
- Select the current location (which is used as the default for new sound measurements)
- View and manage locations via a dropdown menu

To use the frontend, open `frontend/index.html` or `http://localhost:8080` in your browser. Ensure the backend server is running on `localhost:8080`.

## Validation Rules

- **device_id:** Defaults to "arduino_001" if not provided
- **room_name:** Required (auto-filled with current chosen location if omitted)
- **sound_level:** Required, 0-150 dB
- **threshold:** Required, 0-150 dB
- **measure_time:** Auto-generated if not provided
- **is_alert:** Boolean flag
- **description:** Optional text

## Configuration

### Change Authentication Credentials

Edit `internal/api/middleware/basic_authentication.go`:

```go
func validateUser(username, password string) bool {
    return username == "YOUR_USERNAME" && password == "YOUR_PASSWORD"
}
```

### Database Location

Database file: `cmd/api/production.db`

To reset database: Delete the file and restart the server (a new empty database will be created).

## Development Notes

### Default Values

The system automatically provides defaults for:
- `device_id`: "arduino_001"
- `threshold`: 70.0 dB
- `measure_time`: Current timestamp
- `room_name`: Current chosen location (if not provided)

### Logging

Server logs are written to:
- Console output (stdout)
- `cmd/api/production.log`

## Future Enhancements

- [ ] Daily summary endpoint for educators
- [ ] Room-based filtering queries
- [ ] Environment variable configuration
- [ ] JWT authentication
- [ ] WebSocket support for real-time alerts

## License

Educational project for Intelligent Devices course.

---

**Last Updated:** 02/11/2025
