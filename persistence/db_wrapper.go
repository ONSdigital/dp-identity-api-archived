package persistence

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

var nilID = ""

type Cache interface {
	Set(key string, i schema.Identity, ttl time.Duration) error
	Get(key string) (*schema.Identity, error)
	Delete(key string) (bool, error)
}

type CacheWrapper struct {
	IdentityCache Cache
	Database      DB
}

func (c *CacheWrapper) SaveIdentity(newIdentity schema.Identity) (string, error) {
	id, err := c.Database.SaveIdentity(newIdentity)
	if err != nil {
		return nilID, err
	}

	c.IdentityCache.Set(id, newIdentity, 0)
	return id, nil
}

func (c *CacheWrapper) GetIdentity(email string) (schema.Identity, error) {
	i, err := c.IdentityCache.Get(email)
	if err != nil {
		return schema.NilIdentity, err
	}

	if i != nil {
		return *i, nil
	}

	// otherwise request the identity from the database.
	return c.Database.GetIdentity(email)
}
