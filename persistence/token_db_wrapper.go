package persistence

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

var nilID = ""

type Cache interface {
	StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error
	GetToken(ctx context.Context, token string) (time.Duration, error)
}
type CachedTokenStored struct {
	Cache   Cache
	TokenDB TokenStore
}

func (c *CachedTokenStored) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error {

	// c.Cache.StoreToken() .. etc
	// implemented this sprint - for now we'll just fall through

	return c.TokenDB.StoreToken(ctx, tkn, i, ttl)
}

func (c *CachedTokenStored) GetToken(ctx context.Context, tokenStr string) (time.Duration, error) {

	// c.Cache.GetToken() ... etc
	// implemented this sprint - for now we'll just fall through

	return c.TokenDB.GetToken(ctx, tokenStr)
}
