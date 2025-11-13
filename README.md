# Kindergarten Sound Meter

This system supports an IoT device (Arduino with sound sensor) that measures noise levels in kindergarten rooms. When sound exceeds a threshold, data is sent to the server for logging and analysis.

This repository contains:
- A RESTful API backend for monitoring sound levels in kindergarten rooms, built with Go and SQLite.
- A React-based frontend user interface for managing locations, configuring the noise meter, and viewing results and summaries.

## Features

- RESTful API with full CRUD operations for sound data
- SQLite database for data persistence
- Basic authentication
- Sound level validation (0-150 dB range)
- Automatic timestamping
- Alert flagging when threshold exceeded
- **Location management:** Add, select, and manage room locations via `/api/locations` endpoint
- **React Frontend:** Component-based UI with real-time updates
- **Data visualization:** Interactive charts showing daily and weekly noise patterns
- **Real-time monitoring:** Live noise meter with circular progress indicator
- **Statistics dashboard:** Daily peak, weekly average, and monitoring status cards

## Technology Stack

- **Backend Language:** Go 1.x
- **Database:** SQLite
- **Frontend:** React 18, Axios, Bootstrap 5, Recharts
- **Architecture:** Handler-Service-Repository pattern (backend), Component-based (frontend)
- **Testing:** Go testing framework
- **Authentication:** HTTP Basic Auth

## Installation & Setup

### Prerequisites

- Go 1.19 or higher
- Node.js 16+ and npm
- GCC compiler (MinGW or MSYS2 for Windows)

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Enable CGO (if needed)**
   ```bash
   go env -w CGO_ENABLED=1
   ```

### Frontend Setup

1. **Navigate to React frontend directory**
   ```bash
   cd frontend-react
   ```

2. **Install npm dependencies**
   ```bash
   npm install
   ```

## Running the Application

### Development Mode (Recommended)

Run backend and frontend separately with hot-reload enabled.

#### 1. Start Backend Server

```bash
cd backend/cmd/api
go run main.go
```

Backend runs on **http://localhost:8080**

Server will output:
```
Starting server on :8080...
```

#### 2. Start React Development Server

In a new terminal:

```bash
cd frontend-react
npm start
```

Frontend runs on **http://localhost:3000**

Browser will automatically open to http://localhost:3000

**The React app automatically proxies API requests to the backend on port 8080.**

### Production Mode

For deployment, build React app and serve from Go server:

1. **Build React app:**
   ```bash
   cd frontend-react
   npm run build
   ```

2. **Configure Go server to serve React build** (see deployment documentation)

3. **Run Go server:**
   ```bash
   cd backend/cmd/api
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

The React frontend (`frontend-react/`) provides a comprehensive dashboard for monitoring and analyzing classroom noise levels:

### Main Dashboard Components

1. **Settings Panel**
   - Input field for adding locations
   - Location dropdown selector (choose active monitoring room and delete locations)
   - Adjustable noise threshold slider (saved to localStorage)
   - Visual feedback for current settings

2. **Real-time Noise Meter**
   - Circular progress indicator showing current noise level
   - Color-coded status (Quiet/Moderate/Loud)
   - Dynamic updates every 3 seconds
   - Room name display

3. **Statistics Cards**
   - **Daily Peak:** Highest noise level recorded today
   - **Weekly Average:** Average noise level for the current week
   - **Monitoring Status:** Current active room and system status

4. **Noise Analysis Charts**
   - **Daily Analysis:** Area chart showing hourly noise patterns (8:00-17:00)
     - Average and peak noise levels
     - Date selector for historical data
     - Day-of-week selector
   - **Weekly Analysis:** Bar chart comparing noise levels across 5 weekdays
     - Week navigation (up to 4 weeks back)
     - Comparative view of average and peak levels
   - **Independent room selector:** View analytics for any room without changing the active monitoring location

### Technical Features

- **Real-time Updates:** Components automatically refresh with new data
- **Data Persistence:** 
  - LocalStorage for threshold settings
  - SessionStorage for daily peak tracking
- **Error Handling:** User-friendly error messages and loading states
- **Responsive Design:** Works on desktop and mobile devices
- **Interactive Charts:** Built with Recharts library for smooth animations

### Frontend Structure

```
frontend-react/
├── src/
│   ├── components/
│   │   ├── LocationManager.jsx    # Location CRUD operations (Not in use now)
│   │   ├── SettingsPanel.jsx      # Threshold and location controls
│   │   ├── NoiseMeter.jsx          # Circular noise level display
│   │   ├── StatsCards.jsx          # Statistics dashboard cards
│   │   └── NoiseAnalytics.jsx      # Chart visualization component
│   ├── services/
│   │   └── api.js                  # API service layer (axios)
│   ├── App.js                      # Root component with state management
│   ├── App.css                     # Global styles and component styling
│   └── index.js                    # Entry point
├── public/
│   └── index.html
└── package.json
```

## Legacy Frontend

The original HTML/JavaScript frontend is preserved in the `frontend/` directory for reference. It provides basic location management features using vanilla JavaScript.

To use the legacy frontend, open `frontend/index.html` in a browser or access at `http://localhost:8080` while the backend is running.

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

Edit `backend/internal/api/middleware/basic_authentication.go`:

```go
func validateUser(username, password string) bool {
    return username == "YOUR_USERNAME" && password == "YOUR_PASSWORD"
}
```

Also update in `frontend-react/src/services/api.js`:

```javascript
auth: {
  username: 'YOUR_USERNAME',
  password: 'YOUR_PASSWORD'
}
```

### Database Location

Database file: `backend/cmd/api/production.db`

To reset database: Delete the file and restart the server (a new empty database will be created).

## Testing

### Backend Tests

```bash
# Run all tests
cd backend
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/api/handlers/data
```

### Frontend Development

```bash
cd frontend-react

# Start development server with hot-reload
npm start

# Run tests (if configured)
npm test

# Build for production
npm run build
```

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
- `backend/cmd/api/production.log`

### CORS Configuration

The backend allows cross-origin requests from the React development server. CORS headers are configured in `backend/internal/api/middleware/common.go`.

### API Proxy

During development, the React app (port 3000) proxies API requests to the Go backend (port 8080). This is configured in `frontend-react/package.json`:

```json
"proxy": "http://localhost:8080"
```

### Current Data Source

**Note:** The frontend currently uses **simulated data** for demonstration purposes. The real-time noise meter generates random values, and the charts display sample patterns. To connect real Arduino sensor data:

1. Modify the noise simulation logic in `App.js`
2. Replace simulated values with actual sensor readings from the backend API
3. Implement WebSocket connection for true real-time updates (future enhancement)

## Project Structure

```
API 0.1/
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
├── frontend/                 # Legacy HTML/JS frontend
├── frontend-react/           # React frontend (primary UI)
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── services/         # API client
│   │   └── App.js            # Main application
│   └── package.json
└── README.md
```

## Troubleshooting

**CORS errors:**
- Ensure backend is running on port 8080
- Check CORS configuration in `backend/internal/api/middleware/common.go`
- Verify proxy setting in `frontend-react/package.json`

**API 415 errors:**
- Check that Content-Type headers are set correctly in API requests
- Verify interceptor configuration in `frontend-react/src/services/api.js`

**Charts not displaying:**
- Ensure recharts is installed: `npm install recharts`
- Check browser console for React errors
- Verify data format matches chart component expectations

**localStorage/sessionStorage issues:**
- Clear browser storage: Run in console: `sessionStorage.clear(); localStorage.clear();`
- Check browser privacy settings allow storage

## Future Enhancements

- [ ] WebSocket integration for real-time data streaming
- [ ] Historical data export (CSV/PDF)
- [ ] Alert notifications when threshold exceeded
- [ ] Multi-room comparison view
- [ ] Teacher activity correlation analysis
- [ ] JWT authentication replacement
- [ ] Environment variable configuration
- [ ] Unit tests for React components

## License

Educational project for Intelligent Devices course.

---

**Last Updated:** 12/11/2025