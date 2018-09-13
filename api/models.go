package api

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/pkg/errors"
	"time"
)

//go:generate moq -out generate_mocks.go -pkg api . IdentityService

const (
	getIdentityAction = "getIdentity"
	createIdentityAction = "createIdentity"
	identityURIFormat    = "%s/identity/%s"
	headerContentType    = "content-type"
	mimeTypeJSON         = "application/json"
)

var (
	ErrFailedToReadRequestBody      = errors.New("error while attempting to read request body")
	ErrFailedToUnmarshalRequestBody = errors.New("error while attempting to unmarshal request body")
	ErrRequestBodyNil               = errors.New("error expected request body but was empty")
	ErrFailedToGetIdentity = errors.New("error while attempting to get identity")
	ErrFailedToMarshallIdentity = errors.New("error while attempting to marshal identity model")
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
	Get(ctx context.Context) (*identity.Model, error)
}
