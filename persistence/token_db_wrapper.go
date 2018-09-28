package persistence

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

var nilID = ""

type Cache interface {
	StoreToken(key string, i schema.Identity, ttl time.Duration) error
	GetToken(token string) (time.Duration, error)
}
type CacheWrapper struct {
	TokenCache Cache
	Database   TokenStore
}

func (c *CacheWrapper) StoreToken(key string, i schema.Identity, ttl time.Duration) error {

	// until cache is implemented - just fall through
	return c.TokenCache.StoreToken(key, i, ttl)
}

func (c *CacheWrapper) GetToken(token string) (time.Duration, error) {

	// until cache is implemented - just fall through
	return c.TokenCache.GetToken(token)
}
