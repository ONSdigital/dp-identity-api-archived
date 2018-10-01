// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package persistencetest

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"sync"
)

var (
	lockIdentityStoreMockGetIdentity  sync.RWMutex
	lockIdentityStoreMockSaveIdentity sync.RWMutex
)

// IdentityStoreMock is a mock implementation of IdentityStore.
//
//     func TestSomethingThatUsesIdentityStore(t *testing.T) {
//
//         // make and configure a mocked IdentityStore
//         mockedIdentityStore := &IdentityStoreMock{
//             GetIdentityFunc: func(email string) (schema.Identity, error) {
// 	               panic("TODO: mock out the GetIdentity method")
//             },
//             SaveIdentityFunc: func(newIdentity schema.Identity) (string, error) {
// 	               panic("TODO: mock out the SaveIdentity method")
//             },
//         }
//
//         // TODO: use mockedIdentityStore in code that requires IdentityStore
//         //       and then make assertions.
//
//     }
type IdentityStoreMock struct {
	// GetIdentityFunc mocks the GetIdentity method.
	GetIdentityFunc func(email string) (schema.Identity, error)

	// SaveIdentityFunc mocks the SaveIdentity method.
	SaveIdentityFunc func(newIdentity schema.Identity) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetIdentity holds details about calls to the GetIdentity method.
		GetIdentity []struct {
			// Email is the email argument value.
			Email string
		}
		// SaveIdentity holds details about calls to the SaveIdentity method.
		SaveIdentity []struct {
			// NewIdentity is the newIdentity argument value.
			NewIdentity schema.Identity
		}
	}
}

// GetIdentity calls GetIdentityFunc.
func (mock *IdentityStoreMock) GetIdentity(email string) (schema.Identity, error) {
	if mock.GetIdentityFunc == nil {
		panic("moq: IdentityStoreMock.GetIdentityFunc is nil but IdentityStore.GetIdentity was just called")
	}
	callInfo := struct {
		Email string
	}{
		Email: email,
	}
	lockIdentityStoreMockGetIdentity.Lock()
	mock.calls.GetIdentity = append(mock.calls.GetIdentity, callInfo)
	lockIdentityStoreMockGetIdentity.Unlock()
	return mock.GetIdentityFunc(email)
}

// GetIdentityCalls gets all the calls that were made to GetIdentity.
// Check the length with:
//     len(mockedIdentityStore.GetIdentityCalls())
func (mock *IdentityStoreMock) GetIdentityCalls() []struct {
	Email string
} {
	var calls []struct {
		Email string
	}
	lockIdentityStoreMockGetIdentity.RLock()
	calls = mock.calls.GetIdentity
	lockIdentityStoreMockGetIdentity.RUnlock()
	return calls
}

// SaveIdentity calls SaveIdentityFunc.
func (mock *IdentityStoreMock) SaveIdentity(newIdentity schema.Identity) (string, error) {
	if mock.SaveIdentityFunc == nil {
		panic("moq: IdentityStoreMock.SaveIdentityFunc is nil but IdentityStore.SaveIdentity was just called")
	}
	callInfo := struct {
		NewIdentity schema.Identity
	}{
		NewIdentity: newIdentity,
	}
	lockIdentityStoreMockSaveIdentity.Lock()
	mock.calls.SaveIdentity = append(mock.calls.SaveIdentity, callInfo)
	lockIdentityStoreMockSaveIdentity.Unlock()
	return mock.SaveIdentityFunc(newIdentity)
}

// SaveIdentityCalls gets all the calls that were made to SaveIdentity.
// Check the length with:
//     len(mockedIdentityStore.SaveIdentityCalls())
func (mock *IdentityStoreMock) SaveIdentityCalls() []struct {
	NewIdentity schema.Identity
} {
	var calls []struct {
		NewIdentity schema.Identity
	}
	lockIdentityStoreMockSaveIdentity.RLock()
	calls = mock.calls.SaveIdentity
	lockIdentityStoreMockSaveIdentity.RUnlock()
	return calls
}
