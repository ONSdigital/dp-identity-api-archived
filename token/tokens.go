package token

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

//go:generate moq -out tokentest/generate_mocks.go -pkg tokentest . ExpiryTimeHelper Cache

const (
	nilTTL = 0
)

var (
	// ErrTokenExpired returned when GetTTL is called and the token has expired
	ErrTokenExpired = errors.New("token expired")

	errCacheStoreFailed = errors.New("warning failed to write token to cache")
)

// Cache definition of a Token cache.
type Cache interface {
	StoreToken(ctx context.Context, token string, i schema.Identity, ttl time.Duration) error
	GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, time.Duration, error)
}

// ExpiryTimeHelper provides functions for getting the current time and calculating a token's expiry data.
type ExpiryTimeHelper interface {
	Now() time.Time
	GetExpiry() time.Time
}

// Tokens provides functionality for creating new tokens and getting existing ones.
type Tokens struct {
	TimeHelper ExpiryTimeHelper
	Cache      Cache
	Store      persistence.TokenStore
	MaxTTL     time.Duration
}

// NewToken creates and stores a new token for the provided identity. Returns the generated token and its time to live,
// or an error is unsuccessful
func (t *Tokens) NewToken(ctx context.Context, identity schema.Identity) (token *schema.Token, ttl time.Duration, err error) {
	if token, err = t.newToken(identity); err != nil {
		return
	}

	if err = t.Store.StoreToken(ctx, *token, identity); err != nil {
		token = nil
		return
	}

	if ttl, err = t.GetTTL(token); err != nil {
		token = nil
		return
	}

	if err = t.Cache.StoreToken(ctx, token.ID, identity, ttl); err != nil {
		token = nil
		ttl = 0
		return
	}
	log.InfoCtx(ctx, "successfully generated token for identity", log.Data{"identity_id": identity.ID})
	return
}

// Get return the identity associated with the token (if it exists) and the tokens time to live. Return an error if
// unsuccessful
func (t *Tokens) Get(ctx context.Context, tokenStr string) (identity *schema.Identity, ttl time.Duration, err error) {
	defer func() {
		if err != nil { // if err return nil values
			identity = nil
			ttl = 0
		}
	}()

	identity, ttl, err = t.Cache.GetIdentityByToken(ctx, tokenStr)

	if err != nil || identity != nil {
		return
	}

	var token *schema.Token
	if identity, token, err = t.Store.GetIdentityByToken(ctx, tokenStr); err != nil {
		return
	}

	if ttl, err = t.GetTTL(token); err != nil {
		return
	}

	if err = t.Cache.StoreToken(ctx, tokenStr, *identity, ttl); err != nil {
		// We consider this non critical as the token exists and the user can still use the service.
		// So we log an error to record that it happened, clear the error var and carry on.
		log.ErrorCtx(ctx, errors.Wrap(err, errCacheStoreFailed.Error()), log.Data{"identity_id": identity.ID})
		err = nil
	}

	return
}

func (t *Tokens) newToken(i schema.Identity) (*schema.Token, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &schema.Token{
		ID:          uuid.String(),
		IdentityID:  i.ID,
		CreatedDate: t.TimeHelper.Now(),
		ExpiryDate:  t.TimeHelper.GetExpiry(),
		Deleted:     false,
	}, nil
}

// GetTTL calculates the TTL (time to live) from the configured expiry time. Returns ErrTokenExpired if the token is
// expired.
func (t *Tokens) GetTTL(token *schema.Token) (time.Duration, error) {
	if token == nil {
		return nilTTL, errors.New("token required but was nil")
	}
	now := t.TimeHelper.Now()
	if now.After(token.ExpiryDate) {
		return nilTTL, ErrTokenExpired
	}

	// calculate the time remaining until the expiry time
	remainder := token.ExpiryDate.Sub(now)

	if remainder == 0 {
		return nilTTL, ErrTokenExpired
	}

	if remainder.Seconds() >= t.MaxTTL.Seconds() {
		// more than or equal to max TTL so just return max TTL
		return t.MaxTTL, nil
	}

	// time remaining is less than the time until expiry so just return the remaining time.
	return remainder, nil
}
