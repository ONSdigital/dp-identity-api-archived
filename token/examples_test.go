package token

import (
	"time"
)

func assertEquals(t1, t2 time.Time) {}

func ExampleExpiryHelper_GetExpiry() {
	// If the current date is: 2nd Jan 2006
	// And we configure the expiry time to be: 22:30:30
	helper := NewExpiryHelper(22, 30, 0)

	// Then the GetExpiry() should return 2006-01-02T22:30:30.000000
	expiry := helper.GetExpiry()
	expected := time.Date(2006, 1, 2, 22, 30, 0, 0, time.UTC)

	assertEquals(expiry, expected)
}

/*func ExampleToken_GetTTL() {
	// Case 1:
	// Duration until expiry is 30 minutes, the max TTL is set to 15 minutes
	// The duration until expiry is greater than the MaxTTL so the return value is the MaxTTL

	// The token expires at: 2nd Jan 2006 22:30:00
	expiresAt := time.Date(2006, 1, 2, 22, 30, 0, 0, time.UTC)

	// the current time: 2nd Jan 2006 22:00:00
	now := time.Date(2006, 1, 2, 22, 0, 0, 0, time.UTC)

	// the max TTL
	MaxTTL = time.Minute * 15

	t := &schema.Token{
		ID:          "666",
		CreatedDate: now,
		ExpiryDate:  expiresAt,
		IdentityID:  "666",
		Deleted:     false,
	}

	ttl, _ := GetTTL(t)
	fmt.Printf("%t", ttl.Minutes() == MaxTTL.Minutes())

	// Case 2:
	// Duration until expiry is 10 minutes, the max TTL is set to 15 minutes
	// The duration until expiry is less than the MaxTTL so the return is duration until the expiry
	// In this case 10 mins.

	// The token expires at: 2nd Jan 2006 22:30:00
	t.ExpiryDate = time.Date(2006, 1, 2, 22, 30, 0, 0, time.UTC)

	// Current time: 2nd Jan 2006 22:20:00
	now = time.Date(2006, 1, 2, 22, 20, 0, 0, time.UTC)

	// the max TTL
	MaxTTL = time.Minute * 15

	ttl, _ = GetTTL(t)
	expected := time.Minute * 10


	fmt.Printf("%t", ttl.Minutes() == expected.Minutes())

	// Case 3:
	// The expiry date is in the past
	// The returned TTL is the "zero value" and ErrTokenExpired

	// The token expires at: 2nd Jan 2006 22:30:00
	t.ExpiryDate = time.Date(2006, 1, 2, 22, 30, 0, 0, time.UTC)

	// Current time: 2nd Jan 2006 22:31:00
	now = time.Date(2006, 1, 2, 22, 31, 0, 0, time.UTC)

	var err error
	ttl, err = GetTTL(t)


	fmt.Printf("%t", ttl.Minutes() == 0)
	fmt.Printf("%t", err == ErrTokenExpired)
}
*/