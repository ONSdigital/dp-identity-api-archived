package persistence

import (
	"errors"
	"github.com/ONSdigital/dp-identity-api/schema"
)

//go:generate moq -out persistencetest/generate_mocks.go -pkg persistencetest . DB

var (
	ErrNotFound  = errors.New("not found")
	ErrNonUnique = errors.New("non unique")
)

// DB...
type DB interface {
	SaveIdentity(newIdentity schema.Identity) (string, error)
	GetIdentity(email string) (schema.Identity, error)
}