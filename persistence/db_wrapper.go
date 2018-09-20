package persistence

import (
	"github.com/ONSdigital/dp-identity-api/schema"
)

var nilID = ""

type Cache interface {
	Set(key string, i interface{}) error
	Get(key string, i interface{}) error
	Delete(key string) error
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

	c.IdentityCache.Set(id, newIdentity)
	return id, nil
}

func (c *CacheWrapper) GetIdentity(email string) (schema.Identity, error) {
	// First attempt to get the identity from the cache.
	var i schema.Identity
	err := c.IdentityCache.Get(email, &i)
	if err != nil {
		return schema.NilIdentity, err
	}

	// If a not empty return
	if i == schema.NilIdentity {
		return c.Database.GetIdentity(email)
	}

	// otherwise request the identity from the database.
	return schema.NilIdentity, nil
}
