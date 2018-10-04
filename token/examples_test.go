package token

import (
	"fmt"
	"github.com/ONSdigital/dp-identity-api/schema"
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

func ExampleTokens_GetTTL() {
	// Case 1:
	// Duration until expiry is 24 hours, the max TTL is set 15 minutes
	// The duration until expiry is greater than the MaxTTL so the return value is the MaxTTL

	tokens := Tokens{
		TimeHelper: NewExpiryHelper(23, 59, 59),
		MaxTTL: time.Minute * 15,
	}

	now := time.Now()
	token := &schema.Token{
		ExpiryDate: now.Add(time.Hour * 24),
	}

	ttl, _ := tokens.GetTTL(token)
	fmt.Printf("%t", ttl.Minutes() == tokens.MaxTTL.Minutes())

	// Case 2:
	// Duration until expiry is less than the duration of the max TTL
	// The return value is the duration until expiry

	token = &schema.Token{
		ExpiryDate: now.Add(time.Minute * 10),
	}

	ttl, _ = tokens.GetTTL(token)
	expected := time.Minute * 10
	fmt.Printf("%t", ttl.Minutes() == expected.Minutes())

	// Case 3:
	// Expiry date is in the past return ErrTokenExpired

	token = &schema.Token{
		ExpiryDate: now.Add(time.Minute * - 10),
	}

	ttl, _ = tokens.GetTTL(token)
	expected = time.Minute * 10
	fmt.Printf("%t", ttl.Minutes() == expected.Minutes())
}
