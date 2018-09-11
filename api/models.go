package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"io"
	"io/ioutil"
	"time"
)

//go:generate moq -out generate_mocks.go -pkg api . IdentityService

const (
	createIdentityAction = "createIdentity"
	authenticateAction   = "authenticateUser"
	identityURIFormat    = "%s/identity/%s"
	headerContentType    = "content-type"
	mimeTypeJSON         = "application/json"
)

var (
	ErrFailedToReadRequestBody      = errors.New("error while attempting to read request body")
	ErrFailedToUnmarshalRequestBody = errors.New("error while attempting to unmarshal request body")
	ErrRequestBodyNil               = errors.New("error expected request body but was empty")
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

type AuthenticateRequest struct {
	ID    string `json:"id"`
	Password string `json:"password"`
}

type AuthToken struct {
	Token string `json:"token"`
}

func getAuthenticateRequest(r io.ReadCloser) (*AuthenticateRequest, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, ErrFailedToReadRequestBody
	}

	defer r.Close()

	if len(b) == 0 {
		return nil, ErrRequestBodyNil
	}

	var authReq AuthenticateRequest
	if err := json.Unmarshal(b, &authReq); err != nil {
		return nil, ErrFailedToUnmarshalRequestBody
	}
	return &authReq, nil
}

//IdentityService is a service for creating, updating and deleting Identities.
type IdentityService interface {
	Create(ctx context.Context, i *identity.Model) (string, error)
	Authenticate(ctx context.Context, id string, password string) error
}
