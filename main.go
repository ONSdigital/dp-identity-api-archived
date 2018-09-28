package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/api"
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/encryption"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/dp-identity-api/mongo"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	mongolib "github.com/ONSdigital/go-ns/mongo"
	"github.com/ONSdigital/go-ns/server"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serviceNamespace = "dp-identity-api"

func main() {
	log.Namespace = serviceNamespace

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		log.ErrorC("error loading service configuration", err, nil)
		os.Exit(1)
	}

	mongodb, err := mongo.New(cfg.MongoConfig)
	if err != nil {
		log.ErrorC("failed to initialise mongo, exiting app", err, nil)
		os.Exit(1)
	}

	healthTicker := healthcheck.NewTicker(
		cfg.HealthCheckInterval,
		cfg.HealthCheckTimeout,
		mongolib.NewHealthCheckClient(mongodb.Session),
	)

	// use Nop until kafka is added to environment
	auditor := &audit.NopAuditor{}

	apiErrors := make(chan error, 1)

	identityService := &identity.Service{
		IdentityStore: mongodb,
		TokenStore:    mongodb,
		Encryptor:     encryption.Service{},
	}

	identityAPI := api.New("http://localhost"+cfg.BindAddr, identityService, auditor) // TODO make Host config

	router := mux.NewRouter()
	identityAPI.RegisterEndpoints(router)

	httpServer := startHTTPServer(cfg.BindAddr, router, apiErrors)

	for {
		select {
		case err := <-apiErrors:
			log.ErrorC("api error received shutting down service", err, nil)
			gracefulShutdown(cfg.GracefulShutdownTimeout, httpServer, healthTicker, mongodb.Session)
		case s := <-signals:
			log.Debug("os signal received shutting down service", log.Data{"signal": s.String()})
			gracefulShutdown(cfg.GracefulShutdownTimeout, httpServer, healthTicker, mongodb.Session)
		}
	}
}

//startHTTPServer creates and starts a new HTTP Server for the service.
func startHTTPServer(bindAddr string, router *mux.Router, errorChan chan error) *server.Server {
	httpServer := server.New(bindAddr, router)

	go func() {
		log.Debug("starting identity api http server", nil)
		if err := httpServer.ListenAndServe(); err != nil {
			log.ErrorC("httpServer.ListenAndServe() returned an error", err, nil)
			errorChan <- err
		}
	}()
	return httpServer
}

//gracefulShutdown attempts to gracefully shutdown the service resources before existing.
func gracefulShutdown(timeout time.Duration, httpServer *server.Server, healthTicker *healthcheck.Ticker, mongoSess *mgo.Session) {
	log.Info(fmt.Sprintf("shutdown with timeout: %s", timeout), nil)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// stop any incoming requests before closing any outbound connections
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error(err, nil)
	}

	healthTicker.Close()

	if err := mongolib.Close(ctx, mongoSess); err != nil {
		log.Error(err, nil)
	}

	log.Info("shutdown complete", nil)

	cancel()
	os.Exit(1)
}
