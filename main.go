package main

import (
	"github.com/ONSdigital/dp-identity-api/api"
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/mongo"
	"github.com/ONSdigital/dp-identity-api/store"
	"github.com/ONSdigital/go-ns/log"
	"os"
)

const serviceNamespace = "dp-identity-api"

func main() {

	log.Namespace = serviceNamespace
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

	api.CreateIdentityAPI(*store, *cfg)

}
