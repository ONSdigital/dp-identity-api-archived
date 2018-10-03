// Identity package provides functionality for creating/updating/deleting identities and verifying a password for a
// given identity.
package identity

import (
	"errors"
	"github.com/ONSdigital/dp-identity-api/persistence"
)

//go:generate moq -out identitytest/generate_mocks.go -pkg identitytest . Encryptor

// Identity errors.
var (
	// ErrInvalidArguments the required fields with empty or invalid
	ErrInvalidArguments = errors.New("error while attempting create new identity")

	// ErrPersistence an error occurred while attempting to save/get/update an identity from the data store.
	ErrPersistence = errors.New("error while attempting to write data to mongo")

	// ErrAuthenticateFailed authentication was unsuccessful
	ErrAuthenticateFailed = errors.New("authentication unsuccessful")

	// ErrEmailAlreadyExists email is already associated with an active identity
	ErrEmailAlreadyExists = errors.New("active identity already exists with email")

	// ErrIdentityNotFound the requested identity does not exist.
	ErrIdentityNotFound = errors.New("authentication unsuccessful user not found")
)

// Encryptor is a service for encrypting and comparing passwords.
type Encryptor interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

//Service encapsulates the logic for creating, updating and deleting identities
type Service struct {
	IdentityStore persistence.IdentityStore
	TokenStore    persistence.Cache
	Encryptor     Encryptor
}
