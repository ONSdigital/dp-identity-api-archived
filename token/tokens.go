package token

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
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
)

type Cache interface {
	StoreToken(ctx context.Context, token string, i schema.Identity, ttl time.Duration) error
	GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, time.Duration, error)
}

type ExpiryTimeHelper interface {
	Now() time.Time
	GetExpiry() time.Time
}

type Tokens struct {
	TimeHelper ExpiryTimeHelper
	Cache      Cache
	Store      persistence.TokenStore
	MaxTTL     time.Duration
}

func (t *Tokens) NewToken(ctx context.Context, identity schema.Identity) (*schema.Token, error) {
	token, err := t.newToken(identity)
	if err != nil {
		return nil, err
	}

	// TODO remove TTL param.
	if err := t.Store.StoreToken(ctx, *token, identity); err != nil {
		return nil, err
	}

	ttl, _ := t.GetTTL(token)
	if err := t.Cache.StoreToken(ctx, token.ID, identity, ttl); err != nil {
		return nil, err
	}

	return token, nil
}

func (t *Tokens) Get(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error) {
	// try the cache first..
	identity, ttl, err := t.Cache.GetIdentityByToken(ctx, tokenStr)
	if err != nil {
		return nil, 0, err
	}

	// exists in the cache (lovely jubbly) return
	if identity != nil {
		return identity, ttl, nil
	}

	// else fall through to DB.
	var token *schema.Token
	identity, token, err = t.Store.GetIdentityByToken(ctx, tokenStr)
	if err != nil {
		return nil, 0, err
	}

	// calculate TTL
	ttl, err = t.GetTTL(token)
	if err != nil {
		return nil, 0, err
	}

	// put it in the cache for next time.
	if err := t.Cache.StoreToken(ctx, tokenStr, *identity, ttl); err != nil {
		return nil, 0, err
	}

	// happy days...
	return identity, ttl, nil
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
