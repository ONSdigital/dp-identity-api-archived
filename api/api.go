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
	"net/http"
)

var httpServer *server.Server

type IdentityAPI struct {
	dataStore          store.DataStore
	host               string
	router             *mux.Router
	healthCheckTimeout time.Duration
	auditor            audit.AuditorService
}

type apiError struct {
	status  int
	message string
}

func (err *apiError) Error() string {
	return err.message
}

func New(storer store.Storer, cfg config.Configuration, auditor audit.AuditorService) *IdentityAPI {
	api := &IdentityAPI{
		dataStore:          store.DataStore{Backend: storer},
		host:               "http://localhost:" + cfg.BindAddr,
		router:             mux.NewRouter(),
		healthCheckTimeout: cfg.HealthCheckTimeout,
		auditor:            auditor,
	}

	api.router.HandleFunc("/identity", api.CreateIdentityHandler).Methods("POST")
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

func (api *IdentityAPI) CreateIdentityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if auditErr := api.auditor.Record(ctx, createIdentityAction, audit.Attempted, nil); auditErr != nil {
		http.Error(w, auditErr.Error(), http.StatusInternalServerError)
		return
	}

	apiErr := api.createIdentity(ctx, r)
	if apiErr != nil {
		api.auditor.Record(ctx, createIdentityAction, audit.Unsuccessful, nil)
		http.Error(w, apiErr.message, apiErr.status)
		return
	}

	err := api.auditor.Record(ctx, createIdentityAction, audit.Successful, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.InfoCtx(ctx, "createIdentity: identity created successfully", nil)
}
