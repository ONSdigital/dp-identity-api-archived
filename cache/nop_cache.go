package cachepackage

import (
	"github.com/ONSdigital/go-ns/log"
	"time"
	"github.com/ONSdigital/dp-identity-api/schema"
)

// NOPCache is a no op implementation of a cache.
type NOPCache struct{}

func (c *NOPCache) StoreToken(key string, i schema.Identity, ttl time.Duration) error {

	log.Info("nopcache: store token", log.Data{"key": key})
	return nil
}

func (c *NOPCache) GetToken(token string) (time.Duration, error) {

	log.Info("nopcache: get token", log.Data{"token": token})
	return 0, nil
}
