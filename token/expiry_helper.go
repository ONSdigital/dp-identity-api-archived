package token

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"time"
)

const (
	maxHour    = 23
	maxMin     = 59
	maxSec     = 59
	defaultVal = 0
)

var (
	invalidTimeFMT = "%s invalid, must be gte 0 and lt %d defaulting to 0"
)

// ExpiryHelper provides helper functions for calculating token expiry and TTL times.
type ExpiryHelper struct {
	expiryHour   int
	expiryMinute int
	expirySecond int
}

// NewExpiryHelper construct a new NewExpiryHelper instance. Params expiryHour, expiryMinute, expirySecond specify
// the time to at which newly created tokens will expire.
//
// expiryHour, expiryMinute, expirySecond default to 0 if the provided testCases are outside of the ranges [0-23] (hour)
// [0-59] (minutes/seconds).
func NewExpiryHelper(expiryHour, expiryMinute, expirySecond int) *ExpiryHelper {
	if expiryHour <= 0 || expiryHour > maxHour {
		expiryHour = defaultVal
		log.Info(fmt.Sprintf(invalidTimeFMT, "expiryHr", maxHour), nil)
	}

	if expiryMinute <= 0 || expiryMinute > maxMin {
		expiryMinute = defaultVal
		log.Info(fmt.Sprintf(invalidTimeFMT, "expiryMin", maxMin), nil)
	}

	if expirySecond <= 0 || expirySecond > maxSec {
		expirySecond = defaultVal
		log.Info(fmt.Sprintf(invalidTimeFMT, "expirySec", maxSec), nil)
	}

	helper := &ExpiryHelper{
		expiryHour:   expiryHour,
		expiryMinute: expiryMinute,
		expirySecond: expirySecond,
	}

	expiry := helper.GetExpiry()
	log.Info(fmt.Sprintf("token expiry time: %s", expiry.Format("15:04:05")), nil)
	return helper
}

// Now return the current time.
func (e *ExpiryHelper) Now() time.Time {
	return time.Now()
}

// GetExpiry calculate the expiry time for a token relative to the current date.
func (e *ExpiryHelper) GetExpiry() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), e.expiryHour, e.expiryMinute, e.expirySecond, 0, time.UTC)
}
