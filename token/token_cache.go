package token

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

var nilID = ""

type Cache interface {
	StoreToken(ctx context.Context, token schema.Token, i schema.Identity, ttl time.Duration) error
	GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, time.Duration, error)
}

type CachedStore struct {
	Cache Cache
	Store persistence.TokenStore
}

func (c *CachedStore) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error {

	// c.Cache.StoreToken() .. etc
	// implemented this sprint - for now we'll just fall through

	return c.Store.StoreToken(ctx, tkn, i, ttl)
}

func (c *CachedStore) GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {

	// c.Cache.GetToken() ... etc
	// implemented this sprint - for now we'll just fall through

	return c.Store.GetIdentityByToken(ctx, token)
}
