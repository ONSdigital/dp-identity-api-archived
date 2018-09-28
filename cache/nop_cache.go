package cache

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"time"
)

// NOPCache is a no op implementation of a cache.
type NOPCache struct{}


// SetToken is a stand along method for creating a token
func (c *NOPCache) SetToken(key string, i schema.Identity, ttl time.Duration) error {
	log.Info("nopcache: set identity", log.Data{
		"key":   key,
		"ID":    i.ID,
		"email": i.Email,
		"ttl":   ttl.Seconds(),
	})
	return nil
}

// StoreToken .... uses the cache where it can
func (c *NOPCache) StoreToken(token string, identityID string) (time.Duration, error) {
	log.Info("nopcache: store token", log.Data{"key": token})
	return 0, nil
}

// GetToken .... uses the cache where it can
func (c *NOPCache) GetToken(token string, identityID string) (time.Duration, error) {
	log.Info("nopcache: get token", log.Data{"token": token, "identityID":identityID})
	return 0, nil
}

// GetIdentity is implemented in the cache interface for fall-through, but we do not directly serve identities via the cache
func (c *NOPCache) GetIdentity(key string) (*schema.Identity, error) {return nil, nil}

// DeleteIdentity is implemented in the cache interface for fall-through, but we do not directly delete identities via the cache
func (c *NOPCache) DeleteIdentity(key string) (bool, error) {return false, nil}

// CreateIdentity is implemented in the cache interface for fall-through, but we do not directly create identities via the cache
func (c *NOPCache) SaveIdentity(newIdentity schema.Identity) (string, error) {return "", nil}


