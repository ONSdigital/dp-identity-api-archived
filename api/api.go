package api

import (
	"github.com/gorilla/mux"
	"time"

	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/mongo"
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

func CreateIdentityAPI(mongodb *mongo.Mongo, cfg config.Configuration, errorChan chan error) {

	router := mux.NewRouter()

	var auditor audit.AuditorService
	auditor = &audit.NopAuditor{}

	store := store.DataStore{Backend:mongodb}

	api := &IdentityAPI{
		dataStore:          store,
		host:               "http://localhost:" + cfg.BindAddr,
		router:             router,
		healthCheckTimeout: cfg.HealthCheckTimeout,
		auditor:            auditor,
	}

	api.router.HandleFunc("/identity", api.CreateIdentity).Methods("POST")

	// 'Real' routes
	api.router.Path("/healthcheck").HandlerFunc(healthcheck.Do)

	httpServer = server.New(cfg.BindAddr, router)

	go func() {
		log.Debug("Starting api...", nil)
		if err := httpServer.ListenAndServe(); err != nil {
			log.ErrorC("api http server returned error", err, nil)
			errorChan <- err
		}
	}()

}

// Close represents the graceful shutting down of the http server
func Close(ctx context.Context) error {
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}
	log.Info("graceful shutdown of http server complete", nil)
	return nil
}