package identity

import (
	"errors"
	"github.com/ONSdigital/dp-identity-api/mongo"
)

//go:generate moq -out generate_mocks.go -pkg identity . Persistence Encryptor

var (
	ErrInvalidArguments = errors.New("error while attempting create new identity")
	ErrPersistence      = errors.New("error while attempting to write data to mongo")
)

// Persistence...
type Persistence interface {
	Create(newIdentity mongo.Identity) (string, error)
	GetIdentity(id string) (*mongo.Identity, error)
}

type Encryptor interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type ValidationErr struct {
	message string
}

func (e ValidationErr) Error() string {
	return e.message
}

//Service encapsulates the logic for creating, updating and deleting identities
type Service struct {
	Persistence Persistence
	Encryptor   Encryptor
}

//Model is an object representation of a user identity.
type Model struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	UserType          string `json:"user_type"`
	TemporaryPassword bool   `json:"temporary_password"`
	Migrated          bool   `json:"migrated"`
	Deleted           bool   `json:"deleted"`
}
