package api

import (
	"github.com/gorilla/mux"
	"time"

	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/store"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"context"
)

var httpServer *server.Server

type IdentityAPI struct {
	dataStore          store.DataStore
	host               string
	router             *mux.Router
	healthCheckTimeout time.Duration
	auditor            audit.AuditorService
}

func New(storer store.Storer, cfg config.Configuration, auditor audit.AuditorService) *IdentityAPI {
	api := &IdentityAPI{
		dataStore:          store.DataStore{Backend: storer},
		host:               "http://localhost:" + cfg.BindAddr,
		router:             mux.NewRouter(),
		healthCheckTimeout: cfg.HealthCheckTimeout,
		auditor:            auditor,
	}

	api.router.HandleFunc("/identity", api.CreateIdentity).Methods("POST")
	api.router.Path("/healthcheck").HandlerFunc(healthcheck.Do)

	return api
}

func (api *IdentityAPI) GetRouter() *mux.Router {
	return api.router
}

// Close represents the graceful shutting down of the http server
func Close(ctx context.Context) error {
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}
	log.Info("graceful shutdown of http server complete", nil)
	return nil
}
