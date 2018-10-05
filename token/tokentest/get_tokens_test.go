package tokentest

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence/persistencetest"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/dp-identity-api/token"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var (
	testID = "666"
)

func TestTokens_GetCacheError(t *testing.T) {
	Convey("given cache.GetIdentityByToken returns an error", t, func() {
		cache := &CacheMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
				return nil, 0, errTest
			},
		}

		store := &persistencetest.TokenStoreMock{StoreTokenFunc: dbStoreTokenNoErr}

		tokens := token.Tokens{
			Cache:  cache,
			Store:  store,
			MaxTTL: testTTL,
		}

		Convey("when get token is called", func() {
			identity, ttl, err := tokens.GetIdentityByToken(context.Background(), testID)

			Convey("then expected identity, TTL and error values are returned", func() {
				So(err, ShouldEqual, errTest)
				So(identity, ShouldBeNil)
				So(ttl, ShouldEqual, 0)

				So(cache.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(cache.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)
				So(store.GetIdentityByTokenCalls(), ShouldHaveLength, 0)
			})
		})
	})
}

func TestTokens_GetExistsInCache(t *testing.T) {
	Convey("given the request token exists in the cache", t, func() {
		cache := &CacheMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
				return testIdentity, time.Minute * 15, nil
			},
		}

		store := &persistencetest.TokenStoreMock{StoreTokenFunc: dbStoreTokenNoErr}

		tokens := token.Tokens{
			Cache:  cache,
			Store:  store,
			MaxTTL: testTTL,
		}

		Convey("when get token is called", func() {
			identity, ttl, err := tokens.GetIdentityByToken(context.Background(), testID)

			Convey("then expected identity, TTL and error values are returned", func() {
				So(identity, ShouldResemble, testIdentity)
				So(ttl, ShouldEqual, time.Minute*15)
				So(err, ShouldBeNil)

				So(cache.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(cache.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)
				So(store.GetIdentityByTokenCalls(), ShouldHaveLength, 0)
			})
		})
	})
}

func TestTokens_GetStoreError(t *testing.T) {
	Convey("given store.GetIdentityByToken returns an error", t, func() {
		cache := &CacheMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
				return nil, 0, nil
			},
		}

		store := &persistencetest.TokenStoreMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, *schema.Token, error) {
				return nil, nil, errTest
			},
		}

		tokens := token.Tokens{
			Cache:  cache,
			Store:  store,
			MaxTTL: testTTL,
		}

		Convey("when get token is called", func() {
			identity, ttl, err := tokens.GetIdentityByToken(context.Background(), testID)

			Convey("then expected identity, TTL and error values are returned", func() {
				So(err, ShouldEqual, errTest)
				So(identity, ShouldBeNil)
				So(ttl, ShouldEqual, 0)

				So(cache.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(cache.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)

				So(store.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(store.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)
			})
		})
	})
}

func TestTokens_GetStoreTokenExpired(t *testing.T) {
	Convey("given store.GetIdentityByToken returns an expired token", t, func() {
		created := time.Now().Add(time.Hour * - 24) // created a day ago
		expires := time.Now().Add(time.Hour * -1)   // expired an hour ago
		tkn := newTestToken(created, expires)

		cache := &CacheMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
				return nil, 0, nil
			},
		}

		store := &persistencetest.TokenStoreMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, *schema.Token, error) {
				return testIdentity, tkn, nil
			},
		}

		timeHelp := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return time.Now()
			},
		}

		tokens := token.Tokens{
			Cache:      cache,
			Store:      store,
			TimeHelper: timeHelp,
			MaxTTL:     testTTL,
		}

		Convey("when get token is called", func() {
			identity, ttl, err := tokens.GetIdentityByToken(context.Background(), testID)

			Convey("then expected identity, TTL and error values are returned", func() {
				So(err, ShouldEqual, schema.ErrTokenExpired)
				So(identity, ShouldBeNil)
				So(ttl, ShouldEqual, 0)

				So(cache.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(cache.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)

				So(store.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(store.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)
			})
		})
	})
}

func TestTokens_GetStoreCacheStoreTokenError(t *testing.T) {
	Convey("given cache.StoreToken returns an rrror", t, func() {
		created := time.Now()
		expires := time.Now().Add(time.Hour * 1) // expired an hour ago
		tkn := newTestToken(created, expires)

		cache := &CacheMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
				return nil, 0, nil
			},
			StoreTokenFunc: func(ctx context.Context, token string, i schema.Identity, ttl time.Duration) error {
				return errTest
			},
		}

		store := &persistencetest.TokenStoreMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, *schema.Token, error) {
				return testIdentity, tkn, nil
			},
		}

		timeHelp := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return time.Now()
			},
		}

		tokens := token.Tokens{
			Cache:      cache,
			Store:      store,
			TimeHelper: timeHelp,
			MaxTTL:     testTTL,
		}

		Convey("when get token is called", func() {
			identity, ttl, err := tokens.GetIdentityByToken(context.Background(), testID)

			Convey("then expected identity, TTL and error values are returned", func() {
				So(err, ShouldBeNil)
				So(identity, ShouldResemble, testIdentity)
				So(ttl, ShouldEqual, testTTL)

				So(cache.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(cache.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)

				So(store.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(store.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)

				So(cache.StoreTokenCalls(), ShouldHaveLength, 1)
				So(cache.StoreTokenCalls()[0].Token, ShouldEqual, testID)
			})
		})
	})
}

func TestTokens_GetStoreSuccess(t *testing.T) {
	Convey("given store.GetIdentityByToken returns an non expired token", t, func() {
		created := time.Now().Add(time.Hour * - 24)
		expires := time.Now().Add(time.Hour * 1)
		tkn := newTestToken(created, expires)

		cache := &CacheMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, time.Duration, error) {
				return nil, 0, nil
			},
			StoreTokenFunc: cacheStoreTokenNoErr,
		}

		store := &persistencetest.TokenStoreMock{
			GetIdentityByTokenFunc: func(ctx context.Context, token string) (*schema.Identity, *schema.Token, error) {
				return testIdentity, tkn, nil
			},
		}

		timeHelp := &ExpiryTimeHelperMock{
			NowFunc: func() time.Time {
				return time.Now()
			},
		}

		tokens := token.Tokens{
			Cache:      cache,
			Store:      store,
			TimeHelper: timeHelp,
			MaxTTL:     testTTL,
		}

		Convey("when get token is called", func() {
			identity, ttl, err := tokens.GetIdentityByToken(context.Background(), testID)

			Convey("then expected identity, TTL and error values are returned", func() {
				So(err, ShouldBeNil)
				So(identity, ShouldResemble, testIdentity)
				So(ttl, ShouldEqual, testTTL)

				So(cache.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(cache.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)

				So(store.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(store.GetIdentityByTokenCalls()[0].Token, ShouldEqual, testID)

				So(timeHelp.NowCalls(), ShouldHaveLength, 1)
			})
		})
	})
}
