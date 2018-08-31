package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/ONSdigital/go-ns/log"
)

// Configuration structure which hold information for configuring the import API
type Configuration struct {
	BindAddr                string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval     time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckTimeout      time.Duration `envconfig:"HEALTHCHECK_TIMEOUT"`
	MongoConfig             MongoConfig
}

// MongoConfig contains the config required to connect to MongoDB.
type MongoConfig struct {
	BindAddr   string `envconfig:"MONGODB_BIND_ADDR"   json:"-"`
	Collection string `envconfig:"MONGODB_COLLECTION"`
	Database   string `envconfig:"MONGODB_DATABASE"`
}

var cfg *Configuration

// Get the application and returns the configuration structure
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		BindAddr:                ":23800",
		GracefulShutdownTimeout: 5 * time.Second,
		HealthCheckInterval:     30 * time.Second,
		HealthCheckTimeout:      2 * time.Second,
		MongoConfig: MongoConfig{
			BindAddr:   "localhost:27017",
			Collection: "identities",
			Database:   "identities",
		},
	}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	// sensitive fields are omitted from config.String().
	log.Info("loaded service configuration", log.Data{"config": cfg})
	return cfg, nil
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Configuration) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
