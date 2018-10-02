package mongo

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

func (m *Mongo) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error {
	return nil
}

func (m *Mongo) GetToken(ctx context.Context, tokenStr string) (time.Duration, error) {
	return time.Second * 0, nil
}
