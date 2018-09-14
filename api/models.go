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
	getIdentityAction    = "getIdentity"
	createIdentityAction = "createIdentity"
	createToken          = "createToken"
	identityURIFormat    = "%s/identity/%s"
	headerContentType    = "content-type"
	mimeTypeJSON         = "application/json"
)

var (
	ErrFailedToReadRequestBody      = errors.New("error while attempting to read request body")
	ErrFailedToUnmarshalRequestBody = errors.New("error while attempting to unmarshal request body")
	ErrRequestBodyNil               = errors.New("error expected request body but was empty")
	ErrNoTokenProvided = errors.New("error expected token was not provided.")
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

type NewTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getNewTokenRequest(ctx context.Context, r io.ReadCloser) (*NewTokenRequest, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "error reading request body"), nil)
		return nil, ErrFailedToReadRequestBody
	}

	defer r.Close()

	if len(b) == 0 {
		log.ErrorCtx(ctx, errors.Wrap(err, "new token request body expected but was empty"), nil)
		return nil, ErrRequestBodyNil
	}

	var authReq NewTokenRequest
	if err := json.Unmarshal(b, &authReq); err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "error unmarshaling new token request body"), nil)
		return nil, ErrFailedToUnmarshalRequestBody
	}

	if authReq.Email == "" {
		log.ErrorCtx(ctx, errors.New("new token request email expected but was empty"), nil)
		return nil, ErrAuthRequestIDNil
	}
	return &authReq, nil
}

//IdentityService is a service for creating, updating and deleting Identities.
type IdentityService interface {
	Create(ctx context.Context, i *identity.Model) (string, error)
	Get(ctx context.Context) (*identity.Model, error)
	CreateToken(ctx context.Context, email string, password string) error
}
