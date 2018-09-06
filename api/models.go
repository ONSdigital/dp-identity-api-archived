package api

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"net/http"
	"time"
)

//go:generate moq -out generate_mocks.go -pkg api . IdentityService

const (
	createIdentityAction = "createIdentity"
	identityURIFormat    = "%s/identity/%s"
	headerContentType    = "content-type"
	mimeTypeJSON         = "application/json"
)

var (
	ErrFailedToReadRequestBody      = errors.New("error while attempting to read request body")
	ErrFailedToUnmarshalRequestBody = errors.New("error while attempting to unmarshal request body")
	ErrRequestBodyNil               = errors.New("error expected request body but was empty")

	//map specific errors to http status codes.
	errorStatusMapping = map[error]int{
		ErrFailedToUnmarshalRequestBody: http.StatusInternalServerError,
		ErrFailedToReadRequestBody:      http.StatusInternalServerError,
		ErrRequestBodyNil:               http.StatusBadRequest,
		identity.ErrInvalidArguments:    http.StatusInternalServerError,
		identity.ErrPersistence:         http.StatusInternalServerError,
		identity.ErrNameValidation:      http.StatusBadRequest,
		identity.ErrEmailValidation:     http.StatusBadRequest,
		identity.ErrPasswordValidation:  http.StatusBadRequest,
		identity.ErrIdentityNil:         http.StatusBadRequest,
	}
)

//API defines HTTP HandlerFunc's for the endpoints offered by the Identity API service.
type API struct {
	Host               string
	IdentityService    IdentityService
	healthCheckTimeout time.Duration
	auditor            audit.AuditorService
}

//IdentityCreated is the HTTP response entity for create identity success.
type IdentityCreated struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

//IdentityService is a service for creating, updating and deleting Identities.
type IdentityService interface {
	Create(ctx context.Context, i *identity.Model) (string, error)
}
