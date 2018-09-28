package token

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const testIdentityID = "666"

func TestNew(t *testing.T) {
	Convey("should create new token with expected values", t, func() {
		now := time.Now()
		currentTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 40, 0, 0, time.UTC)
		expiryTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

		// The Tokens struct has a func value to return the current time. In the real world we would simple pass in
		// a func which returns time.Now(). For unit testing we pass in a func which returns a time value we've
		// specified which enables us to accurately verify the calculated "TTL".
		getCurrentTimeFunc := func() time.Time {
			return currentTime
		}

		getExpiryFunc := func() time.Time {
			return expiryTime
		}

		Init(maxTTL,  getCurrentTimeFunc, getExpiryFunc)
		token, err := New(testIdentityID)

		So(err, ShouldBeNil)
		So(token.ID, ShouldNotBeEmpty)
		So(token.IdentityID, ShouldEqual, testIdentityID)
		So(token.CreatedDate, ShouldEqual, currentTime)
		So(token.ExpiryDate, ShouldEqual, expiryTime)
		So(token.Deleted, ShouldBeFalse)
	})


}

func TestTokens_GetTTLShouldReturnMaxTTL(t *testing.T) {
	Convey("should return max TTL if the time remaining until expiry is greater than the max TTL", t, func() {
		maxTTL := time.Minute * 15
		now := time.Now()
		currentTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 40, 0, 0, time.UTC)
		expiryTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

		// The Tokens struct has a func value to return the current time. In the real world we would simple pass in
		// a func which returns time.Now(). For unit testing we pass in a func which returns a time value we've
		// specified which enables us to accurately verify the calculated "TTL".
		getCurrentTimeFunc := func() time.Time {
			return currentTime
		}

		getExpiryFunc := func() time.Time {
			return expiryTime
		}

		Init(maxTTL,  getCurrentTimeFunc, getExpiryFunc)

		tokens, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := tokens.GetTTL()

		So(err, ShouldBeNil)
		So(ttl.Seconds(), ShouldEqual, maxTTL.Seconds())
	})
}
func TestTokens_GetTTLShouldReturnTimeRemaining(t *testing.T) {
	Convey("should return TTL equal to the time remaining if the time until expiry is less than the max TTL", t, func() {
		maxTTL := time.Minute * 15
		now := time.Now()
		// 5 mins until expiry
		currentTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 54, 59, 0, time.UTC)
		expiryTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

		// The Tokens struct has a func value to return the current time. In the real world we would simple pass in
		// a func which returns time.Now(). For unit testing we pass in a func which returns a time value we've
		// specified which enables us to accurately verify the calculated "TTL".
		getCurrentTimeFunc := func() time.Time {
			return currentTime
		}

		getExpiryFunc := func() time.Time {
			return expiryTime
		}

		Init(maxTTL,  getCurrentTimeFunc, getExpiryFunc)

		tokens, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := tokens.GetTTL()

		expected := time.Minute * 5
		So(err, ShouldBeNil)
		So(ttl.Seconds(), ShouldEqual, expected.Seconds())
	})
}

func TestTokens_GetTTLShouldReturnExpired(t *testing.T) {
	Convey("should return token expired error if current time equals the expiry time", t, func() {
		maxTTL := time.Minute * 15
		now := time.Now()
		expiryTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 00, 00, 0, time.UTC)
		currentTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 00, 00, 0, time.UTC)

		getCurrentTimeFunc := func() time.Time {
			return currentTime
		}

		getExpiryFunc := func() time.Time {
			return expiryTime
		}

		Init(maxTTL,  getCurrentTimeFunc, getExpiryFunc)

		tokens, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := tokens.GetTTL()
		So(err, ShouldEqual, ErrTokenExpired)
		So(ttl, ShouldEqual, 0)
	})

	Convey("should return token expired error if current time is after the expiry time", t, func() {
		maxTTL := time.Minute * 15
		now := time.Now()
		expiryTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 00, 00, 0, time.UTC)
		currentTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 01, 00, 0, time.UTC)

		getCurrentTimeFunc := func() time.Time {
			return currentTime
		}

		getExpiryFunc := func() time.Time {
			return expiryTime
		}

		Init(maxTTL,  getCurrentTimeFunc, getExpiryFunc)

		tokens, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := tokens.GetTTL()
		So(err, ShouldEqual, ErrTokenExpired)
		So(ttl, ShouldEqual, 0)
	})
}
