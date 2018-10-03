package cache

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

// NOP is a no op implementation of a cache.
type NOP struct{}

func (c *NOP) StoreToken(ctx context.Context, token string, i schema.Identity, ttl time.Duration) error {
	return nil
}

func (c *NOP) GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
	return nil, 0, nil
}