package server

import (
	"context"
	"goapi/internal/api/handlers/data"
	"goapi/internal/api/handlers/locations"
	"goapi/internal/api/middleware"
	"goapi/internal/api/service"
	"log"
	"net/http"
	"path/filepath"
)

type Server struct {
	ctx        context.Context
	HTTPServer *http.Server
	logger     *log.Logger
}

func NewServer(ctx context.Context, sf *service.ServiceFactory, logger *log.Logger) *Server {
	// Create a separate mux for API to apply authentication middleware
	apiMux := http.NewServeMux()
	if err := setupDataHandlers(apiMux, sf, logger); err != nil {
		logger.Fatalf("Error setting up data handlers: %v", err)
	}
	if err := setupLocationHandlers(apiMux, sf, logger); err != nil {
		logger.Fatalf("Error setting up location handlers: %v", err)
	}

	// Main mux serves frontend static files and mounts API under /api/
	mux := http.NewServeMux()
	frontendDir := filepath.Join("..", "..", "..", "frontend")
	absFrontendDir, _ := filepath.Abs(frontendDir)
	logger.Println("Serving frontend from:", absFrontendDir)
	mux.Handle("/", http.FileServer(http.Dir(absFrontendDir)))

	// Apply authentication & common middleware to API
	middlewares := []middleware.Middleware{
		middleware.BasicAuthenticationMiddleware,
		middleware.CommonMiddleware,
	}
	mux.Handle("/api/", http.StripPrefix("/api", middleware.ChainMiddleware(apiMux, middlewares...)))

	return &Server{
		ctx:    ctx,
		logger: logger,
		HTTPServer: &http.Server{
			Handler: mux,
		},
	}
}

func (api *Server) Shutdown() error {
	api.logger.Println("Gracefully shutting down server...")
	return api.HTTPServer.Shutdown(api.ctx)
}

func (api *Server) ListenAndServe(addr string) error {
	api.HTTPServer.Addr = addr
	return api.HTTPServer.ListenAndServe()
}

// Setup REST API handlers for /data
func setupDataHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {
	ds, err := sf.CreateDataService(service.SQLiteDataService)
	if err != nil {
		return err
	}

	// Standard /data endpoint
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data.GetHandler(w, r, logger, ds)
		case http.MethodPost:
			data.PostHandler(w, r, logger, ds)
		case http.MethodPut:
			data.PutHandler(w, r, logger, ds)
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /data/{id} endpoint
	mux.HandleFunc("/data/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data.GetByIDHandler(w, r, logger, ds)
		case http.MethodDelete:
			data.DeleteHandler(w, r, logger, ds)
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /data/daily/{room} endpoint
	mux.HandleFunc("/data/daily/{room}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data.GetDailySummaryHandler(w, r, logger, ds)
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// ===== NEW: /data/hourly endpoint for Arduino 1-hour average =====
	mux.HandleFunc("/data/hourly", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			data.PostHourlyHandler(w, r, logger, ds) // <-- new hourly handler
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return nil
}

// Setup REST API handlers for /locations
func setupLocationHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {
	ls, err := sf.CreateLocationService()
	if err != nil {
		return err
	}

	mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			locations.GetLocationsHandler(w, r, logger, ls)
		case http.MethodPost:
			locations.CreateLocationHandler(w, r, logger, ls)
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/locations/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			locations.UpdateLocationHandler(w, r, logger, ls)
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			locations.DeleteHandler(w, r, logger, ls)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return nil
}
