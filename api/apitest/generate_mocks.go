// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package apitest

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"sync"
	"time"
)

var (
	lockIdentityServiceMockCreate         sync.RWMutex
	lockIdentityServiceMockVerifyPassword sync.RWMutex
)

// IdentityServiceMock is a mock implementation of IdentityService.
//
//     func TestSomethingThatUsesIdentityService(t *testing.T) {
//
//         // make and configure a mocked IdentityService
//         mockedIdentityService := &IdentityServiceMock{
//             CreateFunc: func(ctx context.Context, i *schema.Identity) (string, error) {
// 	               panic("TODO: mock out the Create method")
//             },
//             VerifyPasswordFunc: func(ctx context.Context, email string, password string) (*schema.Identity, error) {
// 	               panic("TODO: mock out the VerifyPassword method")
//             },
//         }
//
//         // TODO: use mockedIdentityService in code that requires IdentityService
//         //       and then make assertions.
//
//     }
type IdentityServiceMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, i *schema.Identity) (string, error)

	// VerifyPasswordFunc mocks the VerifyPassword method.
	VerifyPasswordFunc func(ctx context.Context, email string, password string) (*schema.Identity, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// I is the i argument value.
			I *schema.Identity
		}
		// VerifyPassword holds details about calls to the VerifyPassword method.
		VerifyPassword []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Email is the email argument value.
			Email string
			// Password is the password argument value.
			Password string
		}
	}
}

// Create calls CreateFunc.
func (mock *IdentityServiceMock) Create(ctx context.Context, i *schema.Identity) (string, error) {
	if mock.CreateFunc == nil {
		panic("moq: IdentityServiceMock.CreateFunc is nil but IdentityService.Create was just called")
	}
	callInfo := struct {
		Ctx context.Context
		I   *schema.Identity
	}{
		Ctx: ctx,
		I:   i,
	}
	lockIdentityServiceMockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	lockIdentityServiceMockCreate.Unlock()
	return mock.CreateFunc(ctx, i)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//     len(mockedIdentityService.CreateCalls())
func (mock *IdentityServiceMock) CreateCalls() []struct {
	Ctx context.Context
	I   *schema.Identity
} {
	var calls []struct {
		Ctx context.Context
		I   *schema.Identity
	}
	lockIdentityServiceMockCreate.RLock()
	calls = mock.calls.Create
	lockIdentityServiceMockCreate.RUnlock()
	return calls
}

// VerifyPassword calls VerifyPasswordFunc.
func (mock *IdentityServiceMock) VerifyPassword(ctx context.Context, email string, password string) (*schema.Identity, error) {
	if mock.VerifyPasswordFunc == nil {
		panic("moq: IdentityServiceMock.VerifyPasswordFunc is nil but IdentityService.VerifyPassword was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Email    string
		Password string
	}{
		Ctx:      ctx,
		Email:    email,
		Password: password,
	}
	lockIdentityServiceMockVerifyPassword.Lock()
	mock.calls.VerifyPassword = append(mock.calls.VerifyPassword, callInfo)
	lockIdentityServiceMockVerifyPassword.Unlock()
	return mock.VerifyPasswordFunc(ctx, email, password)
}

// VerifyPasswordCalls gets all the calls that were made to VerifyPassword.
// Check the length with:
//     len(mockedIdentityService.VerifyPasswordCalls())
func (mock *IdentityServiceMock) VerifyPasswordCalls() []struct {
	Ctx      context.Context
	Email    string
	Password string
} {
	var calls []struct {
		Ctx      context.Context
		Email    string
		Password string
	}
	lockIdentityServiceMockVerifyPassword.RLock()
	calls = mock.calls.VerifyPassword
	lockIdentityServiceMockVerifyPassword.RUnlock()
	return calls
}

var (
	lockTokenServiceMockGetIdentityByToken sync.RWMutex
	lockTokenServiceMockNewToken           sync.RWMutex
)

// TokenServiceMock is a mock implementation of TokenService.
//
//     func TestSomethingThatUsesTokenService(t *testing.T) {
//
//         // make and configure a mocked TokenService
//         mockedTokenService := &TokenServiceMock{
//             GetIdentityByTokenFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error) {
// 	               panic("TODO: mock out the GetIdentityByToken method")
//             },
//             NewTokenFunc: func(ctx context.Context, identity schema.Identity) (*schema.Token, time.Duration, error) {
// 	               panic("TODO: mock out the NewToken method")
//             },
//         }
//
//         // TODO: use mockedTokenService in code that requires TokenService
//         //       and then make assertions.
//
//     }
type TokenServiceMock struct {
	// GetIdentityByTokenFunc mocks the GetIdentityByToken method.
	GetIdentityByTokenFunc func(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error)

	// NewTokenFunc mocks the NewToken method.
	NewTokenFunc func(ctx context.Context, identity schema.Identity) (*schema.Token, time.Duration, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetIdentityByToken holds details about calls to the GetIdentityByToken method.
		GetIdentityByToken []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TokenStr is the tokenStr argument value.
			TokenStr string
		}
		// NewToken holds details about calls to the NewToken method.
		NewToken []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Identity is the identity argument value.
			Identity schema.Identity
		}
	}
}

// GetIdentityByToken calls GetIdentityByTokenFunc.
func (mock *TokenServiceMock) GetIdentityByToken(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error) {
	if mock.GetIdentityByTokenFunc == nil {
		panic("moq: TokenServiceMock.GetIdentityByTokenFunc is nil but TokenService.GetIdentityByToken was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		TokenStr string
	}{
		Ctx:      ctx,
		TokenStr: tokenStr,
	}
	lockTokenServiceMockGetIdentityByToken.Lock()
	mock.calls.GetIdentityByToken = append(mock.calls.GetIdentityByToken, callInfo)
	lockTokenServiceMockGetIdentityByToken.Unlock()
	return mock.GetIdentityByTokenFunc(ctx, tokenStr)
}

// GetIdentityByTokenCalls gets all the calls that were made to GetIdentityByToken.
// Check the length with:
//     len(mockedTokenService.GetIdentityByTokenCalls())
func (mock *TokenServiceMock) GetIdentityByTokenCalls() []struct {
	Ctx      context.Context
	TokenStr string
} {
	var calls []struct {
		Ctx      context.Context
		TokenStr string
	}
	lockTokenServiceMockGetIdentityByToken.RLock()
	calls = mock.calls.GetIdentityByToken
	lockTokenServiceMockGetIdentityByToken.RUnlock()
	return calls
}

// NewToken calls NewTokenFunc.
func (mock *TokenServiceMock) NewToken(ctx context.Context, identity schema.Identity) (*schema.Token, time.Duration, error) {
	if mock.NewTokenFunc == nil {
		panic("moq: TokenServiceMock.NewTokenFunc is nil but TokenService.NewToken was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Identity schema.Identity
	}{
		Ctx:      ctx,
		Identity: identity,
	}
	lockTokenServiceMockNewToken.Lock()
	mock.calls.NewToken = append(mock.calls.NewToken, callInfo)
	lockTokenServiceMockNewToken.Unlock()
	return mock.NewTokenFunc(ctx, identity)
}

// NewTokenCalls gets all the calls that were made to NewToken.
// Check the length with:
//     len(mockedTokenService.NewTokenCalls())
func (mock *TokenServiceMock) NewTokenCalls() []struct {
	Ctx      context.Context
	Identity schema.Identity
} {
	var calls []struct {
		Ctx      context.Context
		Identity schema.Identity
	}
	lockTokenServiceMockNewToken.RLock()
	calls = mock.calls.NewToken
	lockTokenServiceMockNewToken.RUnlock()
	return calls
}
