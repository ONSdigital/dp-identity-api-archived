package persistence

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

var nilID = ""

type Cache interface {
	StoreToken(tkn schema.Token, i schema.Identity, ttl time.Duration) error
	GetToken(token string) (time.Duration, error)
}
type CachedTokenStored struct {
	Cache   Cache
	TokenDB TokenStore
}

func (c *CachedTokenStored) StoreToken(tkn schema.Token, i schema.Identity, ttl time.Duration) error {

	// c.Cache.StoreToken() .. etc
	// implemented this sprint - for now we'll just fall through

	return c.TokenDB.StoreToken(tkn, i, ttl)
}

func (c *CachedTokenStored) GetToken(tokenStr string) (time.Duration, error) {

	// c.Cache.GetToken() ... etc
	// implemented this sprint - for now we'll just fall through

	return c.TokenDB.GetToken(tokenStr)
}
