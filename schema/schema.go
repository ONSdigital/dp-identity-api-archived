package schema

import "time"

var (
	ErrIdentityNil        = ValidationErr{message: "identity required but was nil"}
	ErrNameValidation     = ValidationErr{message: "mandatory field name was empty"}
	ErrEmailValidation    = ValidationErr{message: "mandatory field email was empty"}
	ErrPasswordValidation = ValidationErr{message: "mandatory field password was empty"}
	NilIdentity           = Identity{}
)

type ValidationErr struct {
	message string
}

func (e ValidationErr) Error() string {
	return e.message
}

//Identity is an object representation of a user identity.
type Identity struct {
	ID                string    `bson:"id" json:"id"`
	Name              string    `bson:"name" json:"name"`
	Email             string    `bson:"email" json:"email"`
	Password          string    `bson:"password" json:"password"`
	UserType          string    `bson:"user_type" json:"user_type"`
	TemporaryPassword bool      `bson:"temporary_password" json:"temporary_password"`
	Migrated          bool      `bson:"migrated" json:"migrated"`
	Deleted           bool      `bson:"deleted" json:"deleted"`
	CreatedDate       time.Time `bson:"createdDate" json:"createdDate"`
}

func (i *Identity) Validate() (err error) {
	if i == nil {
		return ErrIdentityNil
	}
	if i.Name == "" {
		return ErrNameValidation
	}
	if i.Email == "" {
		return ErrEmailValidation
	}
	if i.Password == "" {
		return ErrPasswordValidation
	}
	return nil
}
