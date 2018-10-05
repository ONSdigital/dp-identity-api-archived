package tokentest

import (
	"fmt"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/dp-identity-api/token"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTokens_GetTTLTokenExpire(t *testing.T) {
	Convey("should return ErrTokenExpired if token is expired", t, func() {
		now := time.Now()
		tkn := newTestToken(now.Add(time.Hour*- 2), now)

		timeHelper := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return now
			},
		}

		tokens := token.Tokens{TimeHelper: timeHelper, MaxTTL: testTTL}

		ttl, err := tokens.GetTokenTTL(tkn)

		So(err, ShouldEqual, schema.ErrTokenExpired)
		So(ttl, ShouldEqual, 0)
		So(timeHelper.NowCalls(), ShouldHaveLength, 1)
	})
}

func TestTokens_GetTTLTokenNil(t *testing.T) {
	Convey("should return expected error if token is nil", t, func() {
		timeHelper := &ExpiryTimeHelperMock{}

		tokens := token.Tokens{TimeHelper: timeHelper, MaxTTL: testTTL}

		ttl, err := tokens.GetTokenTTL(nil)

		So(err, ShouldEqual, token.ErrTokenNil)
		So(ttl, ShouldEqual, 0)
		So(timeHelper.NowCalls(), ShouldHaveLength, 0)
	})
}

func TestTokens_GetTTLExpiresNow(t *testing.T) {
	Convey("should return expected error if time.now equals expiry time", t, func() {
		now := time.Now()
		tkn := newTestToken(now, now)

		timeHelper := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return now
			},
		}

		tokens := token.Tokens{TimeHelper: timeHelper, MaxTTL: testTTL}

		after := now.After(now)
		fmt.Println(after)

		ttl, err := tokens.GetTokenTTL(tkn)

		So(err, ShouldEqual, schema.ErrTokenExpired)
		So(ttl, ShouldEqual, 0)
		So(timeHelper.NowCalls(), ShouldHaveLength, 1)
	})
}

func TestTokens_GetTTL_RemainderGreaterThanMaxTTL(t *testing.T) {
	Convey("should return TTL is duration until expiry is greater than TTL", t, func() {
		now := time.Now()
		tkn := newTestToken(now, now.Add(time.Hour*1))

		timeHelper := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return now
			},
		}

		tokens := token.Tokens{TimeHelper: timeHelper, MaxTTL: testTTL}

		after := now.After(now)
		fmt.Println(after)

		ttl, err := tokens.GetTokenTTL(tkn)

		So(err, ShouldBeNil)
		So(ttl, ShouldEqual, testTTL)
		So(timeHelper.NowCalls(), ShouldHaveLength, 1)
	})
}

func TestTokens_GetTTL_RemainderLessThanTLL(t *testing.T) {
	Convey("should return remaining duration until expiry if less than TTL", t, func() {
		now := time.Now()
		tkn := newTestToken(now, now.Add(time.Minute*10)) // expires in 10 mins, ttl is 15 mins

		timeHelper := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return now
			},
		}

		tokens := token.Tokens{TimeHelper: timeHelper, MaxTTL: testTTL}

		after := now.After(now)
		fmt.Println(after)

		ttl, err := tokens.GetTokenTTL(tkn)

		So(err, ShouldBeNil)
		So(ttl, ShouldEqual, time.Minute*10)
		So(timeHelper.NowCalls(), ShouldHaveLength, 1)
	})
}

func newTestToken(created, expired time.Time) *schema.Token {
	return &schema.Token{
		ID:           testID,
		IdentityID:   testID,
		CreatedDate:  created,
		LastModified: created,
		ExpiryDate:   expired,
		Deleted:      false,
	}
}
