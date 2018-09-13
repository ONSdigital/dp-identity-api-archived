package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
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

type AuthToken struct {
	Token string `json:"token"`
}

type AuthenticateRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func getAuthenticateRequest(ctx context.Context, r io.ReadCloser) (*AuthenticateRequest, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "error reading request body"), nil)
		return nil, ErrFailedToReadRequestBody
	}

	defer r.Close()

	if len(b) == 0 {
		log.ErrorCtx(ctx, errors.Wrap(err, "authentication request body expected but was empty"), nil)
		return nil, ErrRequestBodyNil
	}

	var authReq AuthenticateRequest
	if err := json.Unmarshal(b, &authReq); err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "error unmarshaling authentication request body"), nil)
		return nil, ErrFailedToUnmarshalRequestBody
	}

	if authReq.ID == "" {
		log.ErrorCtx(ctx, errors.New("authentication request id expected but was empty"), nil)
		return nil, ErrAuthRequestIDNil
	}
	return &authReq, nil
}

//IdentityService is a service for creating, updating and deleting Identities.
type IdentityService interface {
	Create(ctx context.Context, i *identity.Model) (string, error)
	Authenticate(ctx context.Context, id string, password string) error
}
