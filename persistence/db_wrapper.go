package persistence

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

var nilID = ""

type Cache interface {
	SetToken(key string, i schema.Identity, ttl time.Duration) error
	StoreToken(token string, identityID string) (time.Duration, error)
	GetToken(token string, identityID string) (time.Duration, error)
	SaveIdentity(identity schema.Identity) (string, error)
	GetIdentity(key string) (*schema.Identity, error)
	DeleteIdentity(key string) (bool, error)
}

type CacheWrapper struct {
	IdentityCache Cache
	Database      DB
}

func (c *CacheWrapper) SetToken(key string, i schema.Identity, ttl time.Duration) error {

	// TODO - where should this live? Will there ever be fall through to database layer?

	// for now - use cache level method directly.
	return c.IdentityCache.SetToken(key, i, ttl)
}

func (c *CacheWrapper) StoreToken(token string, identityID string) (time.Duration, error) {

	// TODO - part of current "store token" task.
	return c.Database.StoreToken(token, identityID)
}

func (c *CacheWrapper) GetToken(token string, identityID string) (time.Duration, error) {

	// TODO - part of current "get token" task
	return c.Database.GetToken(token, identityID)
}

func (c *CacheWrapper) SaveIdentity(newIdentity schema.Identity) (string, error) {

	// always fall through to the database
	return c.Database.SaveIdentity(newIdentity)

}

func (c *CacheWrapper) GetIdentity(email string) (schema.Identity, error) {

	// always fall through to the database
	return c.Database.GetIdentity(email)
}

func (c *CacheWrapper) DeleteIdentity(key string) (bool, error) {

	// TODO - not implemented
	return false, nil
}
