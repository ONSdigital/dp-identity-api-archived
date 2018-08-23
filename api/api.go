package api

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"time"

	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/store"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"os"
	"github.com/ONSdigital/go-ns/healthcheck"
)

var httpServer *server.Server

type IdentityAPI struct {
	dataStore          store.DataStore
	host               string
	router             *mux.Router
	healthCheckTimeout time.Duration
	auditor            audit.AuditorService
}

func CreateIdentityAPI(store store.DataStore, cfg config.Configuration) {

	router := mux.NewRouter()

	var auditor audit.AuditorService
	var auditProducer kafka.Producer

	auditProducer, err := kafka.NewProducer(cfg.KafkaAddr, cfg.AuditEventsTopic, 0)
	if err != nil {
		log.Error(errors.Wrap(err, "error creating kakfa audit producer"), nil)
		os.Exit(1)
	}

	auditor = audit.New(auditProducer, "dp-identity-api")

	api := &IdentityAPI{
		dataStore:          store,
		host:               "http://localhost:20111",
		router:             router,
		healthCheckTimeout: cfg.HealthCheckTimeout,
		auditor:            auditor,
	}

	// TODO - temporary routes for testing/dev
	api.router.HandleFunc("/identity/{id}", api.GetIdentityByID).Methods("GET")
	api.router.HandleFunc("/identity", api.PostIdentity).Methods("POST")

	// 'Real' routes
	api.router.Path("/healthcheck").HandlerFunc(healthcheck.Do)

	httpServer = server.New(cfg.BindAddr, router)

	log.Debug("Starting api...", nil)
	httpServer.ListenAndServe() // TODO, properly + graceful shutdown

}
