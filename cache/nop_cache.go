package cache

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"time"
)

// NOPCache is a no op implementation of a cache.
type NOPCache struct{}

func (c *NOPCache) Set(key string, i schema.Identity, ttl time.Duration) error {
	log.Info("nopcache: set", log.Data{
		"key":   key,
		"ID":    i.ID,
		"email": i.Email,
		"ttl":   ttl.Seconds(),
	})
	return nil
}

func (c *NOPCache) Get(key string) (*schema.Identity, error) {
	log.Info("nopcache: get", log.Data{"key": key})
	return nil, nil
}

func (c *NOPCache) Delete(key string) (bool, error) {
	log.Info("nopcache: delete", log.Data{"key": key})
	return true, nil
}
