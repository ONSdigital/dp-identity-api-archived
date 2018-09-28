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
	TokenCache    Cache
	Database      TokenStore
}



