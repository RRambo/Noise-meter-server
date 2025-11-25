# Kindergarten Sound Meter

This system supports an IoT device (Arduino with sound sensor) that measures noise levels in kindergarten rooms. When sound exceeds a threshold, data is sent to the server for logging and analysis.

This repository contains:
- A RESTful API backend for monitoring sound levels in kindergarten rooms, built with Go and SQLite.
- A React-based frontend user interface for managing locations, configuring the noise meter, viewing results, summaries, and alert notifications.

## Features

### Backend
- RESTful API with CRUD operations for sound data and locations
- SQLite database with automatic schema creation
- Basic HTTP authentication
- Sound level validation (0-150 dB)
- CORS enabled for React development

### Frontend
- **Real-time Monitoring**: Circular meter with color-coded status (Quiet/Moderate/Loud)
- **Alert System**: Toast notifications + alert history with 3-minute cooldown
- **Statistics Dashboard**: Daily peak, weekly average, current monitoring status
- **Analytics Charts**: Daily/weekly noise patterns with independent room selection
- **Location Management**: Add/delete rooms, set active monitoring location
- **Settings**: Adjustable threshold with localStorage persistence

## Technology Stack

- **Backend**: Go 1.19+, SQLite, Handler-Service-Repository pattern
- **Frontend**: React 18, Axios, Bootstrap 5, Recharts
- **Auth**: HTTP Basic Auth

## Installation & Setup

### Prerequisites

- Go 1.19 with CGO enabled
- Node.js 16+ and npm
- GCC compiler (MinGW or MSYS2 for Windows)

### Installation

```bash
# Backend
cd backend
go mod download
go env -w CGO_ENABLED=1  # if needed

# Frontend
cd frontend-react
npm install
```

## Running the Application

### Development Mode

Run backend and frontend separately with hot-reload enabled.

**Terminal 1 - Backend:**
```bash
cd backend/cmd/api
go run main.go
# Server runs on http://localhost:8080
```

**Terminal 2 - Frontend:**
```bash
cd frontend-react
npm start
# App opens at http://localhost:3000
# API requests auto-proxy to port 8080
```

### Production Mode

For deployment, build React app and serve from Go server:

```bash
cd frontend-react
npm run build
cd ../../backend/cmd/api
go run main.go
```

Access at **http://localhost:8080**

Note that this is specifically for production and deployment. When developing, we are constantly updating the files within `frontend-react` and so couldn't keep running `npm run build`.
That's why we are using access at **http://localhost:3000** in development now.

## API Endpoints

### Authentication

All endpoints require Basic Authentication:
- **Username:** `kids_noisemeter_admin`
- **Password:** `passwordkids`

### Endpoints

| Method | Endpoint                | Description                        |
|--------|-------------------------|------------------------------------|
| GET    | `/api/data`             | Get all sound measurements         |
| GET    | `/api/data/{id}`        | Get specific measurement           |
| GET    | `/api/data/daily/{room}`| Get data based on roomName and time|
| POST   | `/api/data`             | Create new measurement             |
| PUT    | `/api/data`             | Update existing measurement        |
| DELETE | `/api/data/{id}`        | Delete measurement                 |
| GET    | `/api/locations`        | Get all locations                  |
| GET    | `/api/locations/chosen` | Get currently chosen location      |
| POST   | `/api/locations`        | Add a new location                 |
| PUT    | `/api/locations/{id}`   | Set a location as currently chosen |
| DELETE | `/api/locations/{id}`   | Delete a location                  |

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

### Example Requests

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

## Frontend Features

### Main Components
- **SettingsPanel**: Location selector, threshold slider
- **NoiseMeter**: Circular progress indicator with real-time updates, also has line indicator to see the set threshold level
- **StatsCards**: Daily peak, weekly average, monitoring status
- **NoiseAnalytics**: Interactive charts (daily/weekly analysis)
- **AlertToast**: Pop-up notifications (auto-dismiss after 5s)
- **AlertHistory**: List of today's alerts with count badge

### Data Persistence
- **localStorage**: Threshold settings, last alert date
- **sessionStorage**: Daily peak, alert history (clears at midnight)

## Alert System Details

### Trigger Conditions
- Noise level exceeds threshold
- 3-minute cooldown between alerts (prevents spam)

### Notification Features
- Toast appears in top-right corner
- Displays: room name, noise level, timestamp
- Audio alert (custom MP3 or browser beep fallback)
- History persists until midnight (auto-clears daily)

### Audio Setup (Optional)
Place custom sound file at: `frontend-react/public/sounds/alert.mp3`
Falls back to Web Audio API beep if file missing.

### Frontend Structure

```
frontend-react/src/
├── components/
│   ├── SettingsPanel.jsx
│   ├── NoiseMeter.jsx
│   ├── StatsCards.jsx
│   ├── NoiseAnalytics.jsx
│   ├── AlertToast.jsx
│   └── AlertHistory.jsx
├── services/
│   └── api.js              # Axios client with auth
├── utils/
│   └── audioUtils.js       # Alert sound playback
├── styles/                 # Component-specific CSS
└── App.js                  # Root component
```

## Data Validation Rules and Default Values

- **device_id:** Defaults to "arduino_001" if not provided
- **room_name:** Required (auto-filled with current chosen location if omitted)
- **sound_level:** Required, range 0-150 dB
- **threshold:** Required, default 70 db, range 0-150 dB
- **measure_time:** Auro current timestamp if not provided
- **description:** Optional text

## Authentication Configuration

**Backend** (`backend/internal/api/middleware/basic_authentication.go`):
```go
func validateUser(username, password string) bool {
    return username == "YOUR_USERNAME" && password == "YOUR_PASSWORD"
}
```

**Frontend** (`frontend-react/src/services/api.js`):
```javascript
auth: {
  username: 'YOUR_USERNAME',
  password: 'YOUR_PASSWORD'
}
```

### Database Location
- File: `backend/cmd/api/production.db`
- Reset: Delete file and restart server (auto-recreates empty DB)

## Testing

### Backend Tests

```bash
# Backend tests
cd backend
go test ./...
go test -v ./internal/api/handlers/data

# Frontend development
cd frontend-react
npm start        # Dev server with hot-reload
npm test         # Run tests
npm run build    # Production build
```

## Troubleshooting

Server logs are written to `Console output (stdout)` and `backend/cmd/api/production.log`

**CORS Errors**
- Ensure backend is on port 8080
- Check `backend/internal/api/middleware/common.go`
- Verify proxy in `frontend-react/package.json`

**Storage Issues**
- Clear browser storage: `sessionStorage.clear(); localStorage.clear();`
- Check browser privacy settings

**Charts Not Displaying**
- Verify: `npm install recharts`
- Check browser console for errors

**Port 8080 Already in Use**
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID  /F
```

### Current Data Source

**Note:** The frontend currently uses **simulated data** for demonstration purposes. The real-time noise meter generates random values, and the charts display sample patterns. To connect real Arduino sensor data:

1. Modify the noise simulation logic in `App.js`
2. Replace simulated values with actual sensor readings from the backend API
3. Implement WebSocket connection for true real-time updates (future enhancement)

## Project Structure

```
Noise-meter-server/
├── backend/
│   ├── cmd/api/              # Main application entry
│   ├── internal/api/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middleware/       # Auth & CORS
│   │   ├── repository/       # Data access layer
│   │   ├── server/           # Server setup
│   │   └── service/          # Business logic
│   ├── go.mod
│   └── go.sum
├── frontend-react/           # React frontend
└── README.md
```

## Future Enhancements

- [ ] WebSocket integration for real-time data streaming
- [ ] Historical data export (CSV/PDF)
- [ ] Alert notifications when threshold exceeded
- [ ] Toggle button for alerts/notifications
- [ ] Multi-room comparison view
- [ ] Teacher activity correlation analysis
- [ ] JWT authentication replacement
- [ ] Environment variable configuration
- [ ] Unit tests for React components
- [ ] reload button for refreshing all the displayed data

## License

Educational project for Intelligent Devices course.

---

**Last Updated:** 12/11/2025