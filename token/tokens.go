package token

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

const (
	nilToken = ""
	nilTTL   = 0
)

var (
	// ErrTokenExpired returned when GetTTL is called and the token has expired
	ErrTokenExpired = errors.New("token expired")
)

type getTimeFunc func() time.Time

type Tokens struct {
	maxTTLSeconds time.Duration
	getTimeNow    getTimeFunc
}

// New construct a new Tokens struct. maxTTL is the maximum duration a token should exist for in the cache,
// getCurrentTime is function which returns the current time.
func New(maxTTL time.Duration, getCurrentTime getTimeFunc) *Tokens {
	return &Tokens{
		maxTTLSeconds: maxTTL,
		getTimeNow:    getCurrentTime,
	}
}

// NewToken create a new token.
func (t *Tokens) NewToken() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nilToken, err
	}

	return uuid.String(), nil
}

// GetTTL calculates the TTL (time to live) against the configured expiry time.
//
// IF (expiry_time - current_time) >= max_ttl RETURN max_ttl
//
// IF (expiry_time - current_time) < max_ttl RETURN (expiry_time - current_time)
//
// IF expiry_time is in the past OR expiry_time == current_time RETURN ErrTokenExpired
func (t *Tokens) GetTTL(expiryTime time.Time) (time.Duration, error) {
	if t.getTimeNow().After(expiryTime) {
		return nilTTL, ErrTokenExpired
	}

	// calculate the time remaining until the expiry time
	timeRemaining := expiryTime.Sub(t.getTimeNow())

	if timeRemaining == 0 {
		return nilTTL, ErrTokenExpired
	}

	if timeRemaining.Seconds() >= t.maxTTLSeconds.Seconds() {
		// more than or equal to max TTL so just return max TTL
		return t.maxTTLSeconds, nil
	}

	// time remaining is less than the time until expiry so just return the remaining time.
	return timeRemaining, nil
}
