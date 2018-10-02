package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"time"
)

//go:generate moq -out apitest/generate_mocks.go -pkg apitest . IdentityService

const (
	getIdentityAction    = "getIdentity"
	createIdentityAction = "createIdentity"
	createToken          = "createToken"
	identityURIFormat    = "%s/identity/%s"
	headerContentType    = "content-type"
	mimeTypeJSON         = "application/json"
	tokenHeaderKey       = "token"
)

var (
	ErrFailedToReadRequestBody      = errors.New("error while attempting to read request body")
	ErrFailedToUnmarshalRequestBody = errors.New("error while attempting to unmarshal request body")
	ErrRequestBodyNil               = errors.New("error expected request body but was empty")
	ErrNoTokenProvided              = errors.New("error expected token was not provided.")
)

//API defines HTTP HandlerFunc's for the endpoints offered by the Identity API service.
type API struct {
	Host               string
	IdentityService    IdentityService
	TokenService       TokenService
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
	Get(ctx context.Context, tokenStr string) (*schema.Identity, error)
	Create(ctx context.Context, i *schema.Identity) (string, error)
	VerifyPassword(ctx context.Context, email string, password string) error
}

// TokenService is a service for getting and storing tokens
type TokenService interface {
	StoreToken(ctx context.Context, token schema.Token, i schema.Identity, ttl time.Duration) error
	GetToken(ctx context.Context, token string) (time.Duration, error)
}
