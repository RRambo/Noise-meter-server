# Kindergarten Sound Meter

For now, the code in this repo is a RESTful API backend for monitoring sound levels in kindergarten rooms. Built with Go and SQLite.

This system supports an IoT device (Arduino with sound sensor) that measures noise levels in kindergarten rooms. When sound exceeds a threshold, data is sent to the server for logging and analysis.

### Features

- RESTful API with full CRUD operations
- SQLite database for data persistence
- Basic authentication
- Sound level validation (0-150 dB range)
- Automatic timestamping
- Alert flagging when threshold exceeded

## Technology Stack

- **Language:** Go 1.x
- **Database:** SQLite
- **Architecture:** Handler-Service-Repository pattern
- **Testing:** Go testing framework
- **Authentication:** HTTP Basic Auth

## API Endpoints

### Authentication

All endpoints require Basic Authentication:
- **Username:** `kids_noisemeter_admin`
- **Password:** `passwordkids`

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/data` | Get all sound measurements |
| GET | `/data/{id}` | Get specific measurement |
| POST | `/data` | Create new measurement |
| PUT | `/data` | Update existing measurement |
| DELETE | `/data/{id}` | Delete measurement |

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
curl -X POST http://localhost:8080/data \
  -u kindergarten_admin:sound_monitor_2025 \
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

## Testing

### Run Unit Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/api/handlers/data
```

## Validation Rules

- **device_id:** Required, max 50 characters
- **room_name:** Required
- **sound_level:** Required, 0-150 dB
- **threshold:** Required, 0-150 dB
- **measure_time:** Auto-generated if not provided
- **is_alert:** Boolean flag
- **description:** Optional text

## Configuration

### Change Authentication Credentials

Edit `internal/api/middleware/auth.go`:

```go
func validateUser(username, password string) bool {
    return username == "YOUR_USERNAME" && password == "YOUR_PASSWORD"
}
```

### Database Location

Database file: `cmd/api/production.db`

To reset database: Delete the file and restart the server (new empty database will be created).

## Development Notes

### Default Values

The system automatically provides defaults for:
- `device_id`: "arduino_001"
- `threshold`: 70.0 dB
- `measure_time`: Current timestamp

### Logging

Server logs are written to:
- Console output (stdout)
- `cmd/api/production.log`

## Future Enhancements

For identified future enhancements, search "need" in codebase. They are written in comments.
- [ ] Daily summary endpoint for educators
- [ ] Room-based filtering queries
- [ ] Environment variable configuration
- [ ] JWT authentication
- [ ] WebSocket support for real-time alerts

## License

Educational project for Intelligent Devices course.

---

**Last Updated:** 29/10/2025
