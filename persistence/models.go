package persistence

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

//go:generate moq -out persistencetest/generate_mocks.go -pkg persistencetest . IdentityStore TokenStore

var (
	ErrNotFound  = errors.New("not found")
	ErrNonUnique = errors.New("non unique")
)

// IdentityStore...
type IdentityStore interface {
	SaveIdentity(newIdentity schema.Identity) (string, error)
	GetIdentity(email string) (schema.Identity, error)
}

type TokenStore interface {
	StoreToken(ctx context.Context, token schema.Token, i schema.Identity, ttl time.Duration) error
	GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, time.Duration, error)
}
