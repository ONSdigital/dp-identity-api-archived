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

type ExpiryHelper struct {
	expiryHour   int
	expiryMinute int
	expirySecond int
}

// NewExpiryHelper construct a new NewExpiryHelper instance.
//
// expiryHour - The hour at which a token will expire (0-23)
//
// expiryMin - The minute of the expiryHour which a token will expire (0-59)
//
// expirySec - The second of the expiryMin which a token will expire (0-59)
func NewExpiryHelper(expiryHour, expiryMinute, expirySecond int) *ExpiryHelper {
	if expiryHour <= 0 || expiryHour > maxHour {
		expiryHour = defaultVal
		log.Info(fmt.Sprintf(invalidTimeFMT, "expiryHr", maxHour), nil)
	}

	if expiryMinute <= 0 || expiryMinute > maxMin {
		expiryHour = defaultVal
		log.Info(fmt.Sprintf(invalidTimeFMT, "expiryMin", maxMin), nil)
	}

	if expirySecond <= 0 || expirySecond > maxSec {
		expiryHour = defaultVal
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

func (e *ExpiryHelper) Now() time.Time {
	return time.Now()
}

// GetExpiry calculates the expiry time for a token relative to the current date. Returns the current date at the
// configured Expiry time.
//
// Example
//
// If current date is 2nd Jan 2006 and the configured expiry time is 10:30:30pm then the expiry date will be:
//
// 2006-01-02T22:30:30.000000
//
// If current date/time is 2nd Jan 2006 10:29:29pm then the expiry date is still:
//
// 2006-01-02T22:30:30.000000
func (e *ExpiryHelper) GetExpiry() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), e.expiryHour, e.expiryMinute, e.expirySecond, 0, time.UTC)
}
