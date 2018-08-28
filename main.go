package main

import (
	"github.com/ONSdigital/dp-identity-api/api"
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/mongo"
	"github.com/ONSdigital/dp-identity-api/store"
	"github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	mongolib "github.com/ONSdigital/go-ns/mongo"
	"os"
	"fmt"
	"context"
	"os/signal"
	"syscall"
)

const serviceNamespace = "dp-identity-api"

func main() {

	log.Namespace = serviceNamespace
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)


	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	// sensitive fields are omitted from config.String().
	log.Info("loaded config", log.Data{
		"config": cfg,
	})

	mongodb := &mongo.Mongo{
		Collection: cfg.MongoConfig.Collection,
		Database:   cfg.MongoConfig.Database,
		URI:        cfg.MongoConfig.BindAddr,
	}

	session, err := mongodb.Init()
	if err != nil {
		log.ErrorC("failed to initialise mongo", err, nil)
		os.Exit(1)
	}

	mongodb.Session = session

	log.Debug("listening...", log.Data{
		"bind_address": cfg.BindAddr,
	})

	store := &store.DataStore{Backend: *mongodb}

	healthTicker := healthcheck.NewTicker(
		cfg.HealthCheckInterval,
		cfg.HealthCheckTimeout,
		mongolib.NewHealthCheckClient(mongodb.Session),
	)

	apiErrors := make(chan error, 1)

	api.CreateIdentityAPI(*store, *cfg, apiErrors)


	// Gracefully shutdown the application closing any open resources.
	gracefulShutdown := func() {
		log.Info(fmt.Sprintf("shutdown with timeout: %s", cfg.GracefulShutdownTimeout), nil)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)

		// stop any incoming requests before closing any outbound connections
		api.Close(ctx)

		healthTicker.Close()

		if err = mongolib.Close(ctx, session); err != nil {
			log.Error(err, nil)
		}

		log.Info("shutdown complete", nil)

		cancel()
		os.Exit(1)
	}

	for {
		select {
		case err := <-apiErrors:
			log.ErrorC("api error received", err, nil)
			gracefulShutdown()
		case <-signals:
			log.Debug("os signal received", nil)
			gracefulShutdown()
		}
	}
}
