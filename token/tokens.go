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
	ErrTokenExpired = errors.New("token expired")
)

func New(expiryHour int, maxTTL time.Duration) *Tokens {
	if expiryHour < 0 || expiryHour > 23 {
		// default to 0 (12am)
		expiryHour = 0
	}

	return &Tokens{
		maxTTLSeconds: maxTTL,
		expiryHour:    expiryHour,
	}
}

type Tokens struct {
	maxTTLSeconds time.Duration
	// value between 0-23 representing the hour at which the tokens will expire.
	expiryHour int
}

func (t *Tokens) NewToken() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nilToken, err
	}

	return uuid.String(), nil
}

func (t *Tokens) GetTTL(expiryTime time.Time) (time.Duration, error) {
	if time.Now().After(expiryTime) {
		return nilTTL, ErrTokenExpired
	}

	// calculate the time remaining until the expiry time
	totalRemaining := expiryTime.Sub(time.Now())

	if totalRemaining.Seconds() >= t.maxTTLSeconds.Seconds() {
		// more than or equal to max TTL so just return max TTL
		return t.maxTTLSeconds, nil
	}

	// time remaining is less than the time until expiry so just return the remaining time.
	return totalRemaining, nil
}

func (t *Tokens) getExpiryTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), t.expiryHour, 0, 0, 0, time.UTC)
}
