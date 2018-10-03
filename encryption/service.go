// Encryption provides functionality for encrypting an comparing password values. The Service struct implements
// identity.Encryptor and simply adds a wrapper around the Golang bcrypt library - see
// https://godoc.org/golang.org/x/crypto/bcrypt. This gives us the ability to switch our chosen encryption library if we
// choose to with the smallest footprint possible.
package encryption

import "golang.org/x/crypto/bcrypt"

// Service struct implements identity.Encryptor and simply adds a wrapper around the functions we require from
// the "golang.org/x/crypto/bcrypt" library.
type Service struct {
}

// GenerateFromPassword generate a password. See https://godoc.org/golang.org/x/crypto/bcrypt#GenerateFromPassword for
// details.
func (s Service) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// CompareHashAndPassword generate a password. See https://godoc.org/golang.org/x/crypto/bcrypt#CompareHashAndPassword
// for details.
func (s Service) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
