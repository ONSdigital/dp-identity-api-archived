package token

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

//go:generate moq -out tokentest/generate_mocks.go -pkg tokentest . ExpiryTimeHelper

const (
	nilTTL = 0
)

var (
	// ErrTokenExpired returned when GetTTL is called and the token has expired
	ErrTokenExpired = errors.New("token expired")

	TimeHelper ExpiryTimeHelper
	MaxTTL     time.Duration
)

type ExpiryTimeHelper interface {
	Now() time.Time
	GetExpiry() time.Time
}

// Token is a structure that represents an authentication token for the Identity API
type Token struct {
	ID          string    `bson:"token_id"`
	IdentityID  string    `bson:"identity_id"`
	CreatedDate time.Time `bson:"created_date"`
	ExpiryDate  time.Time `bson:"expiry_date"`
	Deleted     bool      `bson:"deleted"`
}

// NewToken create a new token.
func New(identityID string) (*Token, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &Token{
		ID:          uuid.String(),
		IdentityID:  identityID,
		CreatedDate: TimeHelper.Now(),
		ExpiryDate:  TimeHelper.GetExpiry(),
		Deleted:     false,
	}, nil
}

// GetTTL calculates the TTL (time to live) from the configured expiry time. Returns ErrTokenExpired if the token is
// expired.
func (t *Token) GetTTL() (time.Duration, error) {
	now := TimeHelper.Now()
	if now.After(t.ExpiryDate) {
		return nilTTL, ErrTokenExpired
	}

	// calculate the time remaining until the expiry time
	remainder := t.ExpiryDate.Sub(now)

	if remainder == 0 {
		return nilTTL, ErrTokenExpired
	}

	if remainder.Seconds() >= MaxTTL.Seconds() {
		// more than or equal to max TTL so just return max TTL
		return MaxTTL, nil
	}

	// time remaining is less than the time until expiry so just return the remaining time.
	return remainder, nil
}
