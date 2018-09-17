package persistence

import (
	"errors"
)

//go:generate moq -out persistencetest/generate_mocks.go -pkg persistencetest . DB

var (
	ErrNotFound  = errors.New("not found")
	ErrNonUnique = errors.New("non unique")
)

// DB...
type DB interface {
	SaveIdentity(newIdentity Identity) (string, error)
	GetIdentity(email string) (Identity, error)
}

type Identity struct {
	ID                string    `bson:"id"`
	Name              string    `bson:"name"`
	Email             string    `bson:"email"`
	Password          string    `bson:"password"`
	UserType          string    `bson:"user_type"`
	TemporaryPassword bool      `bson:"temporary_password"`
	Migrated          bool      `bson:"migrated"`
	Deleted           bool      `bson:"deleted"`
}
