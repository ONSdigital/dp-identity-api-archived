// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package api

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/identity"
	"sync"
)

var (
	lockIdentityServiceMockCreate sync.RWMutex
)

// IdentityServiceMock is a mock implementation of IdentityService.
//
//     func TestSomethingThatUsesIdentityService(t *testing.T) {
//
//         // make and configure a mocked IdentityService
//         mockedIdentityService := &IdentityServiceMock{
//             CreateFunc: func(ctx context.Context, i *identity.Model) (string, error) {
// 	               panic("TODO: mock out the Create method")
//             },
//         }
//
//         // TODO: use mockedIdentityService in code that requires IdentityService
//         //       and then make assertions.
//
//     }
type IdentityServiceMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, i *identity.Model) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// I is the i argument value.
			I *identity.Model
		}
	}
}

// Create calls CreateFunc.
func (mock *IdentityServiceMock) Create(ctx context.Context, i *identity.Model) (string, error) {
	if mock.CreateFunc == nil {
		panic("moq: IdentityServiceMock.CreateFunc is nil but IdentityService.Create was just called")
	}
	callInfo := struct {
		Ctx context.Context
		I   *identity.Model
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
	I   *identity.Model
} {
	var calls []struct {
		Ctx context.Context
		I   *identity.Model
	}
	lockIdentityServiceMockCreate.RLock()
	calls = mock.calls.Create
	lockIdentityServiceMockCreate.RUnlock()
	return calls
}
