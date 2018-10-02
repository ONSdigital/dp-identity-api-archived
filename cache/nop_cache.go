package cache

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"time"
)

// NOPCache is a no op implementation of a cache.
type NOPCache struct{}

func (c *NOPCache) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error {

	log.Info("nopcache: store token", log.Data{"key": tkn.ID})

	return nil
}

func (c *NOPCache) GetToken(ctx context.Context, token string) (time.Duration, error) {

	log.Info("nopcache: get token", log.Data{"token": token})

	return 0, nil
}
