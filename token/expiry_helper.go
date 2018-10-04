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
	timeFMT    = "15:04:05" // go time fmt value for HH:MM:SS
)

var (
	invalidTimeFMT = "invalid time value, must be gte 0 and lt %d defaulting to 0"
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
	helper := &ExpiryHelper{
		expiryHour:   getValOrDefault(expiryHour, maxHour),
		expiryMinute: getValOrDefault(expiryMinute, maxMin),
		expirySecond: getValOrDefault(expirySecond, maxSec),
	}

	expiry := helper.GetExpiry()
	log.Info(fmt.Sprintf("token expiry time: %s", expiry.Format(timeFMT)), nil)
	return helper
}

func (e *ExpiryHelper) GetExpiryHour() int {
	return e.expiryHour
}

func (e *ExpiryHelper) GetExpiryMin() int {
	return e.expiryMinute
}

func (e *ExpiryHelper) GetExpirySec() int {
	return e.expirySecond
}

func getValOrDefault(val int, max int) int {
	if val <= 0 || val > max {
		log.Info(fmt.Sprintf(invalidTimeFMT, max), nil)
		return defaultVal
	}
	return val
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
