package identity

import (
	"errors"
	"github.com/ONSdigital/dp-identity-api/persistence"
)

//go:generate moq -out identitytest/generate_mocks.go -pkg identitytest . Encryptor

var (
	ErrInvalidArguments = errors.New("error while attempting create new identity")
	ErrPersistence      = errors.New("error while attempting to write data to mongo")
)

type Encryptor interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

//Service encapsulates the logic for creating, updating and deleting identities
type Service struct {
	DB        persistence.DB
	Encryptor Encryptor
}
