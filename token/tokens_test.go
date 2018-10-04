package token

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence/persistencetest"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/dp-identity-api/token/tokentest"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var (
	testIdentity = &schema.Identity{
		Name:              "John Paul Jones",
		Email:             "blackdog@ons.gov.uk",
		Password:          "foo",
		UserType:          "bar",
		TemporaryPassword: false,
		Migrated:          false,
		Deleted:           false,
		ID:                "666",
		CreatedDate:       time.Now(),
	}

	testTTL = time.Minute * 15

	dbStoreTokenNoErr = func(ctx context.Context, token string, i schema.Identity, ttl time.Duration) error {
		return nil
	}
	cacheStoreTokenNoErr = func(e context.Context, token schema.Token, i schema.Identity) error {
		return nil
	}
)

func TestTokens_NewTokenMaxTTL(t *testing.T) {
	Convey("given new token does not return an error", t, func() {
		now := time.Now()

		cache := &tokentest.CacheMock{StoreTokenFunc: dbStoreTokenNoErr}

		store := &persistencetest.TokenStoreMock{StoreTokenFunc: cacheStoreTokenNoErr}

		timeHelper := &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return now.Add(time.Hour * 24)
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		tokens := Tokens{
			Cache:      cache,
			Store:      store,
			TimeHelper: timeHelper,
			MaxTTL:     testTTL,
		}

		token, ttl, err := tokens.NewToken(context.Background(), *testIdentity)
		So(err, ShouldBeNil)

		Convey("then store.StoreToken should be called 1 time with the expected params", func() {
			So(store.StoreTokenCalls(), ShouldHaveLength, 1)
			So(store.StoreTokenCalls()[0].Token, ShouldResemble, *token)
			So(store.StoreTokenCalls()[0].I, ShouldResemble, *testIdentity)
		})

		Convey("and cache.StoreToken should be called 1 time with the expected params", func() {
			So(cache.StoreTokenCalls(), ShouldHaveLength, 1)
			So(cache.StoreTokenCalls()[0].Token, ShouldEqual, token.ID)
			So(cache.StoreTokenCalls()[0].I, ShouldResemble, *testIdentity)
			So(cache.StoreTokenCalls()[0].TTL, ShouldEqual, testTTL)

		})

		Convey("and the correct TTL is returned", func() {
			So(ttl, ShouldEqual, testTTL)
		})
	})
}

func TestTokens_NewTokenLessThatMaxTTL(t *testing.T) {
	Convey("given duration until expiry is less than the max TTL", t, func() {
		now := time.Now()

		cache := &tokentest.CacheMock{StoreTokenFunc: dbStoreTokenNoErr}

		store := &persistencetest.TokenStoreMock{StoreTokenFunc: cacheStoreTokenNoErr}

		timeHelper := &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return now.Add(time.Minute * 5)
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		tokens := Tokens{
			Cache:      cache,
			Store:      store,
			TimeHelper: timeHelper,
			MaxTTL:     testTTL,
		}

		token, ttl, err := tokens.NewToken(context.Background(), *testIdentity)
		So(err, ShouldBeNil)

		Convey("then store.StoreToken should be called 1 time with the expected params", func() {
			So(store.StoreTokenCalls(), ShouldHaveLength, 1)
			So(store.StoreTokenCalls()[0].Token, ShouldResemble, *token)
			So(store.StoreTokenCalls()[0].I, ShouldResemble, *testIdentity)
		})

		Convey("and cache.StoreToken should be called 1 time with the expected params", func() {
			So(cache.StoreTokenCalls(), ShouldHaveLength, 1)
			So(cache.StoreTokenCalls()[0].Token, ShouldEqual, token.ID)
			So(cache.StoreTokenCalls()[0].I, ShouldResemble, *testIdentity)
			So(cache.StoreTokenCalls()[0].TTL, ShouldEqual, time.Minute*5)

		})

		Convey("and the correct TTL equals the duration until expiry", func() {
			So(ttl, ShouldEqual, time.Minute*5)
		})
	})
}
