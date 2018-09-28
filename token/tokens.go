package token

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

const (
	nilTTL = 0
)

var (
	// ErrTokenExpired returned when GetTTL is called and the token has expired
	ErrTokenExpired = errors.New("token expired")

	getTimeNow    getTimeFunc
	getExpiryDate getTimeFunc
	maxTTL        time.Duration
)

type getTimeFunc func() time.Time

// Token is a structure that represents an authentication token for the Identity API
type Token struct {
	ID          string    `bson:"token_id"`
	IdentityID  string    `bson:"identity_id"`
	CreatedDate time.Time `bson:"created_date"`
	ExpiryDate  time.Time `bson:"expiry_date"`
	Deleted     bool      `bson:"deleted"`
}

// Init initialises the token package.
func Init(maxTokenDuration time.Duration, getTimeNowFunc getTimeFunc, getExpiryDateFunc getTimeFunc) {
	maxTTL = maxTokenDuration
	getTimeNow = getTimeNowFunc
	getExpiryDate = getExpiryDateFunc
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
		CreatedDate: getTimeNow(),
		ExpiryDate:  getExpiryDate(),
		Deleted:     false,
	}, nil
}

// GetTTL calculates the TTL (time to live) against the configured expiry time.
//
// IF (expiry_time - current_time) >= max_ttl RETURN max_ttl
//
// IF (expiry_time - current_time) < max_ttl RETURN (expiry_time - current_time)
//
// IF expiry_time is in the past OR expiry_time == current_time RETURN ErrTokenExpired
func (t *Token) GetTTL() (time.Duration, error) {
	if getTimeNow().After(t.ExpiryDate) {
		return nilTTL, ErrTokenExpired
	}

	// calculate the time remaining until the expiry time
	timeRemaining := t.ExpiryDate.Sub(getTimeNow())

	if timeRemaining == 0 {
		return nilTTL, ErrTokenExpired
	}

	if timeRemaining.Seconds() >= maxTTL.Seconds() {
		// more than or equal to max TTL so just return max TTL
		return maxTTL, nil
	}

	// time remaining is less than the time until expiry so just return the remaining time.
	return timeRemaining, nil
}

// Expire update the tokens deleted flag to true.
func (t *Token) Expire() {
	t.Deleted = true
}
