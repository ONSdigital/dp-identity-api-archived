package identity

import "net/http"

//go:generate moq -out generate_mocks.go -pkg identity . Persistence

var (
	ErrInvalidArguments = &ServiceError{
		status:  http.StatusInternalServerError,
		message: "error while attempting create new identity",
	}

	ErrFailedToReadRequestBody = &ServiceError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to read request body",
	}

	ErrFailedToUnmarshalRequestBody = &ServiceError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to unmarshal request body",
	}

	ErrFailedToWriteToMongo = &ServiceError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to write data to mongo",
	}
)

// Persistence...
type Persistence interface {
	Create(identity *Model) error
}

//Service...
type Service struct {
	Persistence Persistence
}

type ServiceError struct {
	status  int
	message string
}

func (err *ServiceError) Error() string {
	return err.message
}

type Model struct {
	ID                string `bson:"id" json:"id"`
	Name              string `bson:"name" json:"name"`
	Email             string `bson:"email" json:"email"`
	Password          string `bson:"password" json:"password"`
	TemporaryPassword string `bson:"temporary_password" json:"temporary_password"`
	Migrated          bool   `bson:"migrated" json:"migrated"`
	Deleted           bool   `bson:"deleted" json:"deleted"`
	UserType          string `bson:"user_type" json:"user_type"`
}
