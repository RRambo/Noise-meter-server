package server

import (
	"context"
	"goapi/internal/api/handlers/data"
	"goapi/internal/api/handlers/locations"
	"goapi/internal/api/middleware"
	"goapi/internal/api/service"
	"log"
	"net/http"
)

type Server struct {
	ctx        context.Context
	HTTPServer *http.Server
	logger     *log.Logger
}

func NewServer(ctx context.Context, sf *service.ServiceFactory, logger *log.Logger) *Server {

	mux := http.NewServeMux()
	err := setupDataHandlers(mux, sf, logger)
	if err != nil {
		logger.Fatalf("Error setting up data handlers: %v", err)
	}

	middlewares := []middleware.Middleware{
		middleware.CommonMiddleware,
		middleware.BasicAuthenticationMiddleware,
	}

	err = setupLocationHandlers(mux, sf, logger)
	if err != nil {
		logger.Fatalf("Error setting up location handlers: %v", err)
	}

	return &Server{
		ctx:    ctx,
		logger: logger,
		HTTPServer: &http.Server{
			Handler: middleware.ChainMiddleware(mux, middlewares...),
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

// * REST API handlers
func setupDataHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {

	ds, err := sf.CreateDataService(service.SQLiteDataService)
	if err != nil {
		return err
	}

	mux.HandleFunc("OPTIONS /*", func(w http.ResponseWriter, r *http.Request) {
		data.OptionsHandler(w, r)
	})
	mux.HandleFunc("POST /data", func(w http.ResponseWriter, r *http.Request) {
		data.PostHandler(w, r, logger, ds)
	})
	mux.HandleFunc("PUT /data", func(w http.ResponseWriter, r *http.Request) {
		data.PutHandler(w, r, logger, ds)
	})
	mux.HandleFunc("GET /data", func(w http.ResponseWriter, r *http.Request) {
		data.GetHandler(w, r, logger, ds)
	})
	mux.HandleFunc("GET /data/{id}", func(w http.ResponseWriter, r *http.Request) {
		data.GetByIDHandler(w, r, logger, ds)
	})
	mux.HandleFunc("DELETE /data/{id}", func(w http.ResponseWriter, r *http.Request) {
		data.DeleteHandler(w, r, logger, ds)
	})
	return err

}

func setupLocationHandlers(mux *http.ServeMux, sf *service.ServiceFactory, logger *log.Logger) error {
	ls, err := sf.CreateLocationService()
	if err != nil {
		return err
	}

	mux.HandleFunc("GET /locations", func(w http.ResponseWriter, r *http.Request) {
		locations.GetLocationsHandler(w, r, logger, ls)
	})
	mux.HandleFunc("POST /locations", func(w http.ResponseWriter, r *http.Request) {
		locations.CreateLocationHandler(w, r, logger, ls)
	})
	mux.HandleFunc("PUT /locations/{id}", func(w http.ResponseWriter, r *http.Request) {
		locations.SetChosenLocationHandler(w, r, logger, ls)
	})

	return nil
}
