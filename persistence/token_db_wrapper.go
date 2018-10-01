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
	TokenDb    TokenStore
}

func (c *CacheWrapper) StoreToken(key string, i schema.Identity, ttl time.Duration) error {

	// c.TokenCache.StoreToken() .. etc
	// implemented this sprint - for now we'll just fall through

	return c.TokenDb.StoreToken(key, i, ttl)
}

func (c *CacheWrapper) GetToken(token string) (time.Duration, error) {

	// c.TokenCache.GetToken() ... etc
	// implemented this sprint - for now we'll just fall through

	return c.TokenDb.GetToken(token)
}
