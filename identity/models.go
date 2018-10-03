package identity

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

//go:generate moq -out identitytest/generate_mocks.go -pkg identitytest . Encryptor TokenService

var (
	ErrInvalidArguments = errors.New("error while attempting create new identity")
	ErrPersistence      = errors.New("error while attempting to write data to mongo")
)

type Encryptor interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type TokenService interface {
	NewToken(ctx context.Context, identity schema.Identity) (*schema.Token, time.Duration, error)
	Get(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error)
}

//Service encapsulates the logic for creating, updating and deleting identities
type Service struct {
	IdentityStore persistence.IdentityStore
	Tokens        TokenService
	Encryptor     Encryptor
}
