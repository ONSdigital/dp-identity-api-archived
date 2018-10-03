package token

const testIdentityID = "666"

/*func TestNew(t *testing.T) {
	Convey("should create new token with expected values", t, func() {
		now := time.Now()
		expires := now.Add(time.Hour * 1)

		TimeHelper = &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return expires
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		token, err := New(testIdentityID)

		So(err, ShouldBeNil)
		So(token.ID, ShouldNotBeEmpty)
		So(token.IdentityID, ShouldEqual, testIdentityID)
		So(token.CreatedDate, ShouldEqual, now)
		So(token.ExpiryDate, ShouldEqual, expires)
		So(token.Deleted, ShouldBeFalse)
	})

}

func TestTokens_GetTTLShouldReturnMaxTTL(t *testing.T) {
	Convey("should return max TTL if the time until expiry is greater than the max TTL", t, func() {
		now := time.Now()
		expires := now.Add(time.Hour * 1)

		MaxTTL = time.Minute * 15

		TimeHelper = &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return expires
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		token, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := GetTTL(token)

		So(err, ShouldBeNil)
		So(ttl.Seconds(), ShouldEqual, MaxTTL.Seconds())
	})
}
func TestTokens_GetTTLShouldReturnTimeRemaining(t *testing.T) {
	Convey("should return time until expiry if it is less than max TTL", t, func() {
		now := time.Now()
		expires := now.Add(time.Minute * 10)

		MaxTTL = time.Minute * 15
		TimeHelper = &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return expires
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		token, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := GetTTL(token)

		expected := time.Minute * 10
		So(err, ShouldBeNil)
		So(ttl.Seconds(), ShouldEqual, expected.Seconds())
	})
}

func TestTokens_GetTTLShouldReturnExpired(t *testing.T) {
	Convey("should return token expired error if current time equals the expiry time", t, func() {
		now := time.Now()
		expires := now

		MaxTTL = time.Minute * 15
		TimeHelper = &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return expires
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		token, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := GetTTL(token)
		So(err, ShouldEqual, ErrTokenExpired)
		So(ttl, ShouldEqual, 0)
	})

	Convey("should return token expired error if current time is after the expiry time", t, func() {
		expires := time.Now()                   // expires now
		now := time.Now().Add(time.Minute * 15) // add 15 mins to current time.

		MaxTTL = time.Minute * 15
		TimeHelper = &tokentest.ExpiryTimeHelperMock{
			GetExpiryFunc: func() time.Time {
				return expires
			},
			NowFunc: func() time.Time {
				return now
			},
		}

		token, err := New(testIdentityID)
		So(err, ShouldBeNil)

		ttl, err := GetTTL(token)
		So(err, ShouldEqual, ErrTokenExpired)
		So(ttl, ShouldEqual, 0)
	})
}*/
